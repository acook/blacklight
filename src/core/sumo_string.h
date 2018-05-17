#pragma once

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "./datum.h"

typedef byte* sumo;
typedef byte* cursor;
typedef struct {
  uint32_t len;
  uint32_t cap;
} sumo_header;

static inline sumo sumo_new_prealloc(bl_size cap) {
  bl_size real_cap = cap + sizeof(sumo_header);  // make space for header
#ifdef JEMALLOC_C_
  sumo s = aligned_alloc(real_cap, sizeof(sumo_header));
#else
  // The version of Musl that comes with the old Windows build of ELLCC
  // doesn't expose aligned_alloc or memalign so we fall back to malloc
  sumo s = malloc(real_cap);
#endif
  if (!s) return NULL;  // unable to allocate memory
  sumo_header* h = (void*)s;
  h->len = 0;
  h->cap = cap;
  return s;
}

static sumo sumo_new() { return sumo_new_prealloc(0); }

static bl_size sumolen(sumo s) { return ((sumo_header*)s)->len; }

static bl_size sumocap(sumo s) { return ((sumo_header*)s)->cap; }

static bl_size sumo_sizeof(sumo s) { return sumocap(s) + sizeof(sumo_header); }

static sumo sumo_resize(sumo s, bl_size len) {
  // FIXME: use rallocx to do aligned_realloc when jemalloc available
  sumo s2 = realloc(s, len + sizeof(sumo_header));
  if (s2) {
    ((sumo_header*)s2)->cap = len;
    return s2;
  }
  return s;  // unable to cat due to insufficient allocation
}

static sumo sumo_grow(sumo s, bl_size len) {
  if (sumocap(s) <= len) {
    return sumo_resize(s, len);
  }
  return s;  // resize unnecessary because requested is less than actual
}

static void sumocpy_unsafe(sumo dest, sumo src, bl_size len) {
  memcpy(dest + sizeof(sumo_header), src + sizeof(sumo_header), len);
}

static bl_size sumocpy(sumo dest, sumo src, bl_size len) {
  bl_size src_len = sumolen(src);
  if (src_len < len) len = src_len;
  bl_size dest_cap = sumocap(dest);
  if (dest_cap < len) len = dest_cap;
  if (len) sumocpy_unsafe(src, dest, len);
  return len;
}

static cursor sumo_cursor_new(sumo s) { return s + sizeof(sumo_header); }

static cursor sumo_cursor_mv(sumo s, cursor c, bl_size change) {
  if (((c + change) > (sumo_cursor_new(s))) &&
      ((c + change) < (s + sumo_sizeof(s)))) {
    return c + change;
  }
  return c;  // unable to move cursor to desired position
}

static bl_size sumo_cursor_len(sumo s, cursor c) {
  return (sumolen(s) + sizeof(sumo_header)) - (c - s);
}

static bl_size sumo_cursor_pos(sumo s, cursor c) {
  return c - sumo_cursor_new(s);
}

static sumo sumocat(sumo dest, sumo src) {
  bl_size end = sumocap(dest);
  bl_size len = sumolen(src) + sumolen(dest);

  sumo dest2 = sumo_grow(dest, len);
  if (dest2) {
    memcpy(sumo_cursor_new(dest2) + end, sumo_cursor_new(src), sumolen(src));
    ((sumo_header*)dest2)->len = len;
    return dest2;
  }

  return dest;  // unable to cat due to insufficient allocation
}

static sumo sumocat_str(sumo dest, const char* src) {
  bl_size clen = strlen(src);
  bl_size new_len = clen + sumolen(dest);

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
  bl_size len = sumo_cursor_len(s, c);
  char* str = malloc(len + 1);
  memcpy(str, c, len);
  str[len] = 0x00;  // null-terminate string
  return str;
}

static char* sumo_as_cstr(sumo s) {
  sumo_header* h = (void*)s;
  cursor c = sumo_cursor_new(s);
  size_t len = sumolen(s);
  if (!c[len] || !c[len-1]) {
    return c;
  } else {
    if (sumocap(s) == len) {
      s = sumo_grow(s, len + 1 );
      c = sumo_cursor_new(s);
    }
    c[len] = 0x00;
    return c;
  }
}
