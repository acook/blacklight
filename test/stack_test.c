#include <stdio.h>
#include "../src/core/stack.h"
#include "../src/core/text.h"
#include "../src/core/debug.h"
#include "../ext/local/include/criterion/criterion.h"

Test(stack, all) {
  puts(" - init");
  // initialize
  stack s = stack_new();
  stack_reflect(s);

  puts(" - set d1");
  datum *d1;
  s = stack_fit(s, sizeof(*d1));  // ensure there's enough space, resize if not
  d1 = stack_reserve(s, sizeof(*d1));  // reserve space and return pointer to it
  d1->data = 18446744073709551615ull;
  stack_reflect(s);

  puts(" - set d2");
  datum d2 = {.data = 18446744073709551615ull, .t = Number};
  s = stack_push(s, d2);
  stack_reflect(s);

  datum d3;
  d3 = stack_pop(s);
  datum_reflect(d3);
  stack_reflect(s);

  puts(" - NUMBER");
  d2.t = Number;
  d2.data = 16;
  s = stack_push(s, d2);
  stack_reflect(s);
  d3 = stack_pop(s);
  datum_reflect(d3);

  puts(" - SMALL TEXT (3)");
  utf8 u = calloc(4, 1);
  u[0] = 'F';
  u[1] = 'O';
  u[2] = 'O';
  s = new_text_literal(s, u, 3);
  free(u);
  stack_reflect(s);
  d3 = stack_pop(s);
  datum_reflect(d3);

  puts(" - SMALL TEXT (10)");
  u = calloc(11, 1);
  u[0] = 0x6a;
  u[1] = 0x61;
  u[2] = 0x63;
  u[3] = 0x6b;
  u[4] = 0x68;
  u[5] = 0x61;
  u[6] = 0x6d;
  u[7] = 0x6d;
  u[8] = 0x65;
  u[9] = 0x72;
  s = new_text_literal(s, u, 10);
  free(u);
  stack_reflect(s);
  d3 = stack_pop(s);
  datum_reflect(d3);

  puts(" - LARGE TEXT (30)");
  bl_size ts = 30;
  u = calloc(ts + 1, 1);
  bl_size o;
  for (o = 0; o < ts; o++) {
    u[o] = 'x';
  }
  s = new_text_literal(s, u, ts);
  stack_reflect(s);
  d3 = stack_pop(s);
  datum_reflect(d3);
  free(u);

  cr_assert(utf8cmp(text_from_datum(d3), "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx") == 0);

  utf8 u2 = {'w', 'h', 'a', 't'};

  //return 0;
}
