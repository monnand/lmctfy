#include "clmctfy_mock.h"

#include <string.h>
#include <string.h> // to use strdup
#include <stdarg.h>

#include "macros.h"

struct exec_result *result_list = NULL;

static char *last_error_message;

#define MAXLINE 4096

static void lmctfy_mock_set_error_message_fmt(const char *fmt, va_list ap) {
  char buf[MAXLINE];
  if (fmt == NULL) {
	  return;
  }
  memset(buf, 0, MAXLINE);
  vsnprintf(buf, MAXLINE, fmt, ap);
  last_error_message = strdup(buf);
  return;
}

void lmctfy_mock_set_error_message(const char *fmt, ...) {
  va_list ap;
  va_start(ap, fmt);
  lmctfy_mock_set_error_message_fmt(fmt, ap);
  va_end(ap);
  return;
}

const char *lmctfy_mock_get_last_error_message() {
  return last_error_message;
}

void lmctfy_mock_clear_last_error_message() {
  if (last_error_message != NULL) {
    free(last_error_message);
  }
  last_error_message = NULL;
}

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

void lmctfy_mock_assert_expectations() {
  char buf[MAXLINE];
  struct exec_result *iter = NULL;
  *buf = '\0';
  int n = 0;
  if (result_list != NULL) {
    for (iter = result_list; iter->next != NULL; iter = iter->next) {
      strncat(buf, iter->function_name, MAXLINE); 
      strncat(buf, ",", MAXLINE); 
      n++;
    }
    strncat(buf, iter->function_name, MAXLINE); 
    n++;
    lmctfy_mock_set_error_message("The code you are testing needs to make %d more call(s): %s", n, buf);
  }
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

