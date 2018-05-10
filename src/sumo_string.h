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

static sumo sumo_new() {
#ifdef JEMALLOC_C_
  sumo s = aligned_alloc(sizeof(sumo_header), sizeof(sumo_header));
#else
  // The version of Musl that comes with the old Windows build of ELLCC
  // doesn't expose aligned_alloc or memalign so we fall back to malloc
  sumo s = malloc(sizeof(sumo_header));
#endif
  sumo_header *h = (void*)s;
  h->len = 0;
  h->cap = 0;
  return s;
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
  if (sumocap(s) - len) {
    sumo s2 = realloc(s, len + sizeof(sumo_header));
    if (s2) {
      ((sumo_header*)s2)->cap = len;
      return s2;
    }
  }
  return s;
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

static sumo sumocat(sumo dest, sumo src) {
  size_t end = sumo_sizeof(dest);
  size_t len = sumolen(src) + sumolen(dest);
  sumo dest2 = sumo_resize(dest, len);
  if (dest2)  {
    memcpy(dest + end, src + sizeof(sumo_header), sumolen(src));
    ((sumo_header*)dest)->len = len;
    return dest2;
  }
  return dest;
}

static sumo sumocat_str(sumo dest, const char* src) {
  size_t len = strlen(src);
  sumo dest2 = sumo_resize(dest, sumolen(dest) + len);
  if (dest2) {
    memcpy(dest2 + sizeof(sumo_header), src, len);
    dest2[len] = 0x00;
    ((sumo_header*)dest2)->len = len;
    return dest2;
  }
  return dest;
}

static char* sumo_to_cstr(sumo s) {
  size_t len = sumolen(s);
  char* str = calloc(len + 1, 1);
  memcpy(str, s + sizeof(sumo_header), len);
  str[len] = 0x00;
  return str;
}

static cursor sumo_cursor(sumo s) {
  return s + sizeof(sumo_header);
}
