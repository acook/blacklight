#include <stdio.h>
#include "../ext/utf8.h/utf8.h"
#include "./sumo_string.h"

typedef uint64_t bl_uint;
static void printhex(void *ptr, bl_uint len) {
  byte *seq = ptr;

  bl_uint i;
  for (i = 0; i < len; i++) {
    printf("%02X ", seq[i]);
  }

  printf("\n");
}

static inline void usage(void) { puts("usage: redlight filename.bl"); }

static inline void argcheck(int argc, char *argv[]) {
  if (
      argc != 2 ||
      utf8cmp(argv[1], "--help") == 0 ||
      utf8cmp(argv[1], "-h") == 0) {
    usage();
    exit(0);
  }
}

int main(int argc, char *argv[]) {
  argcheck(argc, argv);

  char* data = argv[1];

  sumo s = sumo_new();
  s = sumocat_str(s, data);
  printf("sumo len: %llu\n", sumolen(s));
  printf("sumo cap: %llu\n", sumocap(s));
  puts("sumo hex:");
  printhex(s, sumo_sizeof(s));
  puts("sumo contents:");
  cursor c = sumo_cursor_new(s); // positions cursor at first user location
  fwrite(c, sumolen(s), 1, stdout);
  puts("");

  c = sumo_cursor_mv(s, c, 10);
  printf("sumo cursor len: %llu\n", sumo_cursor_len(s, c));
  puts("sumo contents at index 10:");

  fwrite(c, sumo_cursor_len(s, c), 1, stdout);
  
  puts("\n");

  char* cstr = sumo_to_cstr(s);
  printf("cstr len: %llu\n", strlen(cstr));
  puts("cstr hex:");
  printhex(cstr, sumolen(s));
  puts("cstr contents:");
  puts(cstr);

  free(cstr);
  free(s);
  return 0;
}
