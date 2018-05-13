#pragma once

#include "./bl_types.h"

typedef union {
  struct {
    byte t;   // type
    byte ex;  // extended type info
    byte x;
    byte y;
    byte z;
    byte xx;
    byte yy;
    byte zz;
    union {
      bl_uint data;
      void *ptr;
    };
  };
  byte b[16];
} datum;
typedef union {
  datum d;
  byte b[sizeof(datum)];
} datum_access;
typedef datum *ptrdatum;
