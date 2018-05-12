#include "../src/core/stack.h"
#include "../src/core/debug.h"

int main(void) {
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
  datum d2 = {.data = 18446744073709551615ull};
  s = stack_push(s, d2);
  stack_reflect(s);

  return 0;
}
