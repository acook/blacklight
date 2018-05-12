#pragma once

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>

typedef uint8_t byte;
typedef byte* sumo;
typedef byte* cursor;
typedef struct {
  uint32_t len;
  uint32_t cap;
} sumo_header;

static inline sumo sumo_new_prealloc(size_t cap) {
  size_t real_cap = cap + sizeof(sumo_header);  // make space for header
#ifdef JEMALLOC_C_
  sumo s = aligned_alloc(real_cap, sizeof(sumo_header));
#else
  // The version of Musl that comes with the old Windows build of ELLCC
  // doesn't expose aligned_alloc or memalign so we fall back to malloc
  sumo s = malloc(real_cap);
#endif
  if (!s) return NULL; // unable to allocate memory
  sumo_header *h = (void*)s;
  h->len = 0;
  h->cap = cap;
  return s;
}

static sumo sumo_new() {
  return sumo_new_prealloc(0);
}

static size_t sumolen(sumo s){
  return ((sumo_header*)s)->len;
}

static size_t sumocap(sumo s) {
  return ((sumo_header*)s)->cap;
}

static size_t sumo_sizeof(sumo s) {
  return sumocap(s) + sizeof(sumo_header);
}

static sumo sumo_resize(sumo s, size_t len) {
  // FIXME: use rallocx to do aligned_realloc when jemalloc available
  sumo s2 = realloc(s, len + sizeof(sumo_header));
  if (s2) {
    ((sumo_header*)s2)->cap = len;
    return s2;
  }
  return s;  // unable to cat due to insufficient allocation
}

static sumo sumo_grow(sumo s, size_t len) {
  if (sumocap(s) <= len) {
    return sumo_resize(s, len);
  }
  return s;  // resize unneccessary because requested is less than actual
}

static void sumocpy_unsafe(sumo dest, sumo src, size_t len) {
  memcpy(dest + sizeof(sumo_header), src + sizeof(sumo_header), len);
}

static size_t sumocpy(sumo dest, sumo src, size_t len) {
  size_t src_len = sumolen(src);
  if (src_len < len) len = src_len;
  size_t dest_cap = sumocap(dest);
  if (dest_cap < len) len = dest_cap;
  if (len) sumocpy_unsafe(src, dest, len);
  return len;
}

static cursor sumo_cursor_new(sumo s) {
  return s + sizeof(sumo_header);
}

static cursor sumo_cursor_mv(sumo s, cursor c, size_t change) {
  if (((c + change) > (sumo_cursor_new(s))) && ((c + change) < (s + sumo_sizeof(s)))) {
    return c + change;
  }
  return c; // unable to move cursor to desired position
}

static size_t sumo_cursor_len(sumo s, cursor c) {
  return (sumolen(s) + sizeof(sumo_header)) - (c - s);
}

static size_t sumo_cursor_pos(sumo s, cursor c) {
  return c - sumo_cursor_new(s);
}

static sumo sumocat(sumo dest, sumo src) {
  size_t end = sumocap(dest);
  size_t len = sumolen(src) + sumolen(dest);

  sumo dest2 = sumo_grow(dest, len);
  if (dest2)  {
    memcpy(sumo_cursor_new(dest2) + end, sumo_cursor_new(src), sumolen(src));
    ((sumo_header*)dest2)->len = len;
    return dest2;
  }

  return dest;  // unable to cat due to insufficient allocation
}

static sumo sumocat_str(sumo dest, const char* src) {
  size_t clen = strlen(src);
  size_t new_len = clen + sumolen(dest);

  sumo dest2 = sumo_grow(dest, new_len);
  if (dest2) {
    memcpy(sumo_cursor_new(dest2), src, clen);
    ((sumo_header*)dest2)->len = new_len;
    return dest2;
  }

  return dest;  // unable to cat due to insufficient allocation
}

static char* sumo_to_cstr(sumo s) {
  cursor c = sumo_cursor_new(s);
  size_t len = sumo_cursor_len(s, c);
  char* str = malloc(len + 1);
  memcpy(str, c, len);
  str[len] = 0x00;  // null-terminate string
  return str;
}
