#pragma once

#include "./bl_types.h"

typedef struct {
  byte t;   // type
  byte ex;  // extended type info
  byte y;
  byte z;
  bl_size len;
  union {
    bl_uint data;
    void *ptr;
  };
} datum;
typedef union {
  datum d;
  byte b[16];
} datumarray;
typedef datum *ptrdatum;
