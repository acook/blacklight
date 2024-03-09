#pragma once

#include "./text.h"
#include "./debug.h"

static void datum_reflect(datum d) {
  if (d.t == Number) {
    printf("number value - %lu\n", d.data);
  } else if (d.t == Text || d.t == (Text ^ Ref)) {
    printf("text contents - %s\n", text_from_datum(d));
  } else {
    printf("unknown type %x - ", d.t);
    printhex(&d, sizeof(datum));
  }
}
