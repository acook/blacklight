#include <stdio.h>
#include "../ext/experimental/lib/files.h"
#include "../ext/experimental/lib/sumo.h"
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

  Sumo filename = sumo_new(utf8size(argv[1]));
  utf8cpy(filename, argv[1]);

  int64_t filesize = fsize(filename);
  if (filesize == -1) {
    printf("file not found: %s\n", filename);
    usage();
    exit(1);
  }

  Sumo contents = sumo_new(filesize);

  FILE *file = fopen(filename, "rb");
  fread(contents, filesize, 1, file);
  contents[filesize] = 0x00; // terminate string in case the last line of the file isn't blank

  printf("file size: %d\n", utf8size(contents));
  printf("file contents:\n%s\n", contents);

  Sumo cursor = sumo_cursor(contents, (bl_uint32){.i=100}); // grab a cursor in the middle of contents
  
  printf("file size: %d\n", utf8size(cursor));
  printf("file contents:\n%s\n", cursor);

  return 0;
}
