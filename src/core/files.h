#include <sys/stat.h>

static off_t fsize(const char *filename) {
  struct stat st;

  if (stat(filename, &st) == 0) return st.st_size;

  return -1;
}