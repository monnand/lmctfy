#ifndef MOCK_CLMCTFY_MACROS_H
#define MOCK_CLMCTFY_MACROS_H

#include "clmctfy.h"
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

struct exec_result {
  int error_code;
  char *message;
  void *function_ptr;
  struct exec_result *next;
};

extern struct exec_result *result_list;

#define POP_RESULT

#define MOCK_FUNCTION_BEGIN(NAME, ...) \
    int NAME (struct status *s, __VA_ARGS__) { \
      if (result_list == NULL   \
          || result_list->function_ptr != NAME) { \
        fprintf(stderr, #NAME " should not be called\n"); \
        exit(-1); \
      } \
      struct exec_result *r = result_list;  \
      if (r->error_code == 0) {  \
        s->error_code = 0;  \
        s->message = NULL;  \
      } else {  \
        s->error_code = r->error_code;  \
        s->message = r->message;  \
      } \
      free(r);  \
      result_list = result_list->next;

#define MOCK_FUNCTION_END   \
      return s->error_code; \
    }


#endif // MOCK_CLMCTFY_MACROS_H
