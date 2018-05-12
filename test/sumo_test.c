#include <assert.h>
#include <string.h>
#include "../src/core/debug.h"
#include "../src/core/sumo_string.h"

int main(void) {
  char* data = "0123456789-=qwertyuiop[]asdfghjkl;'zxcvbnm,./";

  sumo s = sumo_new();
  s = sumocat_str(s, data);
  assert(sumolen(s) == 45);
  assert(sumocap(s) == 45);

  // positions cursor at first user location
  cursor c = sumo_cursor_new(s);
  assert(memcmp(c, data, strlen(data)) == 0);

  sumo s2 = sumo_new();
  s2 = sumocat_str(s2, "zzzyyyxxxwww");
  s = sumocat(s, s2);
  s = sumocat(s, s2);

  // MUST reacquire cursor in case sumo has been reallocated/moved
  c = sumo_cursor_new(s);
  assert(sumo_cursor_len(s, c) == 69);
  assert(sumo_cursor_pos(s, c) == 0);
  // make sure prefix hasn't changed, despite cat
  assert(memcmp(c, data, strlen(data)) == 0);

  c = sumo_cursor_mv(s, c, strlen(data));
  assert(sumo_cursor_len(s, c) == 24);
  assert(sumo_cursor_pos(s, c) == 45);
  char* expected = "zzzyyyxxxwwwzzzyyyxxxwww";
  assert(memcmp(c, expected, sumo_cursor_len(s, c)) == 0);

  char* cstr = sumo_to_cstr(s);
  assert(strlen(cstr) == 69);
  expected =
      "0123456789-=qwertyuiop[]asdfghjkl;'zxcvbnm,./zzzyyyxxxwwwzzzyyyxxxwww";
  assert(memcmp(cstr, expected, strlen(cstr)) == 0);

  return 0;
}
