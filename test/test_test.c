#include "../src/core/debug.h"

static char* utest_suite_name;
static char* utest_file_name;

static int utest_successes = 0;
static int utest_failures = 0;
static int utest_test_count = 0;

// typedef void (*utest_func)(void);
// utest_func *utest_tests;
static char* utest_test_names[10];

static void utest_add_test(char* name) {
  utest_test_names[utest_test_count] = name;
  utest_test_count++;
}

static void utest_fail(char* message, int line) {
  printf("test failed on line (%d) with message: %s\n", line, message);
  utest_failures++;
}

static void utest_pass(void) {
  printf("test passed!\n");
  utest_successes++;
}

//#define FAIL(line, message) utest_fail(line, message)
//#define PASS() utest_pass()
//#define UTEST_EXPECT(line, message, expression) if (condition) {} else
//{FAIL((line), (message));
//#define EXPECT(expression) UTEST_EXPECT(__LINE__, "expression evaluated to
// false")

#define UTEST_CAT_HELPER(_A_, _B_) _A_##_B_
#define UTEST_CAT(_A_, _B_) UTEST_CAT_HELPER(_A_, _B_)
#define UTEST_PROTO(Name) void UTEST_CAT(test_, Name)(void)

#define TEST(Name)      \
  UTEST_PROTO(Name);    \
  UTEST_PROTO(Name)

#define RUN(Name) UTEST_CAT(test_, Name)()

TEST(foo) { return; }

static int utest_run_tests(void) {
  for (int o = 0; o <= utest_test_count; o++) {
    // RUN_TEST(utest_test_names[o]);
    printf("test: %s\n", utest_test_names[o]);
  }

  return !utest_failures;
}

int main(void) {
  int EXITSTATUS = 0;

  RUN(foo);

  return EXITSTATUS;
}
