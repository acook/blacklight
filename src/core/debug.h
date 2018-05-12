#pragma once

#include <stdio.h>
#include <string.h>

static void p(char* str, size_t len) {
  printf(" - : ");
  fwrite(str, len, 1, stdout);
  puts("");
}

static char* strdiffchr(char* s1, char* s2) {
  while (*s1 && *s1 == *s2) {
    s1++;
    s2++;
  }
  return (*s1 == *s2) ? NULL : s1;
}

static void printhex(void *ptr, size_t len) {
  char *seq;
  seq = ptr;

  size_t i;
  for (i = 0; i < len; i++) {
    printf("%02X ", seq[i]);
  }

  printf("\n");
}

static void warn(char* message) {
  fwrite(message, strlen(message), 1, stderr);
  fwrite("\n", 1, 1, stderr);
}
