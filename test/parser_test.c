#include <stdio.h>
#include "../src/core/files.h"
#include "../src/core/sumo_string.h"
#include "../ext/utf8.h/utf8.h"
#include "../src/core/bl_types.h"
#include "../src/core/stack.h"
#include "../src/core/bl_math.h"

static inline void usage(void) { puts("usage: redlight filename.bl"); }

static inline sumo argcheck(int argc, char *argv[]) {
  if (
      argc != 2 ||
      utf8cmp(argv[1], "--help") == 0 ||
      utf8cmp(argv[1], "-h") == 0) {
    usage();
    exit(0);
  }

  sumo filename = sumo_new_prealloc(utf8size(argv[1]));
  utf8cpy(filename, argv[1]);

  return filename;
}

static inline stack parse(sumo code) {
  // stack tokens=stack_new();
  // stack_push(tokens, stack_new());
  return stack_new();
}

int main(int argc, char *argv[]) {
  sumo filename = argcheck(argc, argv);
  int64_t filesize = fsize(filename);
  if (filesize == -1) {
    printf("file not found: %s\n", filename);
    usage();
    exit(1);
  }

  sumo contents = sumo_new_prealloc(filesize);

  FILE *file = fopen(filename, "rb");
  fread(contents, filesize, 1, file);
  contents[filesize] = 0x00; // terminate string in case the last line of the file isn't blank, although this isn't enough for security

  printf("file size: %ld\n", filesize);
  printf("utf8 length: %lu\n", utf8size(contents));
  puts("file contents: (first 100 bytes)"); // need to figure out how to print first 100 runes so we don't bisect graphemes
  fwrite(contents, 1, min(500, filesize), stdout); // can't use printf if contents may include 0x00
  puts("");

  //cursor c = sumo_cursor(contents); // grab a cursor in the middle of contents
  //c = sumo_cursor_mv(contents, c, 100)
  //sumo tokens[];??
  //tokens = parse(contents);

  return 0;
}
