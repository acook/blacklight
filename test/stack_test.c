#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef uint32_t bl_size;
typedef uint64_t bl_uint;
typedef uint8_t byte;
typedef struct {
  bl_size ss;  // total available length in bytes (stack size)
  bl_size sp;  // current used length in bytes (stack pointer)
} stack_header;
typedef byte *stack;
typedef union {
  bl_uint d;
  byte b[sizeof(bl_uint)];
} datum;
typedef datum *ptrdatum;

static void printhex(void *ptr, bl_uint len) {
  byte *seq = ptr;

  bl_uint i;
  for (i = 0; i < len; i++) {
    printf("%02X ", seq[i]);
  }

  printf("\n");
}

static inline stack stack_new() {
// The version of Musl that comes with the old Windows build of ELLCC
// doesn't expose aligned_alloc or memalign so we fall back to malloc
#ifdef JEMALLOC_C_
  stack s = aligned_alloc(sizeof(bl_uint), sizeof(stack_header));
#else
  stack s = malloc(sizeof(stack_header));
#endif
  stack_header *h = (void *)s;
  h->sp = h->ss = sizeof(stack_header);
  return s;
}

// expand stack if not large enough for fit
static stack stack_fit(stack s, bl_size fit) {
  stack_header *h = (void *)s;
  bl_size new_size = (h->sp + fit);
  if (new_size > h->ss) {
    stack new_s = realloc(s, new_size);
    if (new_s == NULL) {
      printf("unable to allocate %u bytes\n", new_size);
    } else {
      printf("realloc'd %u bytes\n", new_size);
      s = new_s;
      h = (void *)s;  // reacquire header address in case it moved
    }
    h->ss = new_size;
  }

  return s;
}

// requires that the stack already has enough space
// reserve space of r bytes, move sp forward same
// return pointer to reserved location
static inline void *stack_reserve(stack s, bl_uint r) {
  stack_header *h = (void *)s;
  void *loc = (void *)(s + h->sp);
  h->sp += r;
  return loc;
}

// push datum onto stack
static inline stack stack_push(stack s, datum d)
    __attribute__((warn_unused_result));
static inline stack stack_push(stack s, datum d) {
  s = stack_fit(s, sizeof(datum));
  *(datum *)stack_reserve(s, sizeof(datum)) = d;
  return s;
}

static void stack_reflect(stack s) {
  stack_header *h = (void *)s;
  printf("stack*: %p\n", (void *)s);
  printf("stack size: %u\n", h->ss);
  printf("stack top: %u\n", h->sp);
  printhex(s, h->sp);
}

int main(void) {
  puts(" - init");
  // initialize
  stack s = stack_new();
  stack_reflect(s);

  puts(" - set d1");
  datum *d1;
  s = stack_fit(s, sizeof(*d1));  // ensure there's enough space, resize if not
  d1 = stack_reserve(s, sizeof(*d1));  // reserve space and return pointer to it
  d1->d = 18446744073709551615ull;
  stack_reflect(s);

  puts(" - set d2");
  datum d2 = {.d = 18446744073709551615ull};
  s = stack_push(s, d2);
  stack_reflect(s);

  return 0;
}
