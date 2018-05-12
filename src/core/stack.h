#pragma once

#include <stdlib.h>
#include <stdint.h>
#include "./debug.h"
#include "./datum.h"

typedef struct {
  bl_size ss;  // total available length in bytes (stack size)
  bl_size sp;  // current used length in bytes (stack pointer)
} stack_header;
typedef byte *stack;

static stack stack_new() {
#ifdef JEMALLOC_C_
  stack s = aligned_alloc(sizeof(bl_uint), sizeof(stack_header));
#else
  // The version of Musl that comes with the old Windows build of ELLCC
  // doesn't expose aligned_alloc or memalign so we fall back to malloc
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

// warning: calling stack_push after this will clobber scalar types!
static inline datum stack_pop(stack s) {
  stack_header *h = (void *)s;
  h->sp = h->sp - sizeof(datum);
  datum *loc = (void *)(s + h->sp);
  return *loc;
}

static void stack_reflect(stack s) {
  stack_header *h = (void *)s;
  printf("stack*: %p\n", (void *)s);
  printf("stack size: %u\n", h->ss);
  printf("stack top: %u\n", h->sp);
  printhex(s, h->sp);
}
