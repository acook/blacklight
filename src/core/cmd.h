#pragma once

#include <stdio.h>
#include "../../ext/utf8.h/utf8.h"

static inline void usage(void) { puts("usage: redlight filename.bl"); }

static inline void argcheck(int argc, char *argv[]) {
  if (argc != 2 || utf8cmp(argv[1], "--help") == 0 ||
      utf8cmp(argv[1], "-h") == 0) {
    usage();
    exit(0);
  }
}
