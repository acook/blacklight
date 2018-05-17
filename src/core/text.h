#pragma once

#include "./stack.h"
#include "./opcode.h"

static stack new_text(stack s) {
  datum d;
  d.t = Text;
  return stack_push(s, d);
}

static stack new_text_literal(stack s, utf8 u, bl_size len) {
  datum d;
  if (len < (sizeof(datum) - sizeof(byte) - sizeof(byte))) {
    printf("text len (%d) can fit onto stack\n", len);
    d.t = Text;
    d.ex = (byte)len;  // can store len here since it will always be 6 or less
    bl_size h = sizeof(byte) + sizeof(byte);
    printf("header size %d\n", h);
    bl_size o;
    for (o = 0; o <= len; o++) {
      ((datumarray)d).b[o + h] = u[o];
    }
  } else {
    printf("text len (%d) too large, putting into ref\n", len);
    d.t = Text ^ Ref;
    d.ptr = u;
  }
  return stack_push(s, d);
}

static inline utf8 text_from_datum(datum d) {
  if (d.t == Text) {
    return &((datumarray)d).b[sizeof(opcode) + sizeof(opcode)];
  } else if (d.t == (Text ^ Ref)) {
    return (utf8)d.ptr;
  } else {
    puts("NOT A TEXT");
    return NULL;
  }
}
