#include <stdio.h>
#include "../src/core/files.h"
#include "../src/core/sumo_string.h"
#include "../ext/utf8.h/utf8.h"

static inline void usage(void) { puts("usage: redlight filename.bl"); }

int main(int argc, char *argv[]) {
  if (
      argc != 2 ||
      utf8cmp(argv[1], "--help") == 0 ||
      utf8cmp(argv[1], "-h") == 0) {
    usage();
    return 0;
  }

  sumo filename = sumo_new_prealloc(utf8size(argv[1]));
  utf8cpy(filename, argv[1]);

  int64_t filesize = fsize(filename);
  if (filesize == -1) {
    printf("file not found: %s\n", filename);
    usage();
    exit(1);
  }

  sumo contents = sumo_new_prealloc(filesize);

  FILE *file = fopen(filename, "rb");
  fread(contents, filesize, 1, file);
  contents[filesize] = 0x00; // terminate string in case the last line of the file isn't blank

  printf("file size: %lu\n", utf8size(contents));
  printf("file contents:\n%s\n", contents);

  cursor c = sumo_cursor_new(contents); // grab a cursor in the middle of contents
  c = sumo_cursor_mv(contents, c, 100);

  printf("file size: %lu\n", utf8size(c));
  printf("file contents:\n%s\n", c);

  return 0;
}
