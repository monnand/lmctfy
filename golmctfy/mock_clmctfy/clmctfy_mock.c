#include "clmctfy_mock.h"
#include "macros.h"
#include <string.h>

struct exec_result *result_list = NULL;

void lmctfy_mock_clear_all_expected_calls() {
  struct exec_result *iter = result_list;
  while (iter != NULL) {
    struct exec_result *ptr = iter;
    iter = ptr->next;
    if (ptr->message != NULL) {
      free(ptr->message);
    }
    free(ptr->function_name);
    free(ptr);
  }
  result_list = NULL;
}

void lmctfy_mock_expect_call(const char *fn, int error_code, const char *message) {
  struct exec_result *r = (struct exec_result *)malloc(sizeof(struct exec_result));
  r->function_name = strdup(fn);
  r->error_code = error_code;
  if (error_code != 0) {
    r->message = strdup(message);
  }
  r->next = NULL;
  struct exec_result *iter = NULL;
  if (result_list == NULL) {
    result_list = r;
    return;
  }
  for (iter = result_list; iter->next != NULL; iter = iter->next) {
  }

  iter->next = r;
}

