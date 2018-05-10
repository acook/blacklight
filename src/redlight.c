#include <stdio.h>
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

int main() {
  sumo s = sumo_new();
  s = sumocat_str(s, "foobar");
  printf("sumo len: %llu\n", sumolen(s));
  printf("sumo cap: %llu\n", sumocap(s));
  puts("sumo hex:");
  printhex(s, sumo_sizeof(s));
  puts("sumo contents:");
  fwrite(sumo_cursor(s), sumolen(s), 1, stdout);
  
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
