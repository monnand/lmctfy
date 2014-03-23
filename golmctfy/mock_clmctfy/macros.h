#ifndef MOCK_CLMCTFY_MACROS_H
#define MOCK_CLMCTFY_MACROS_H

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

#include "clmctfy.h"
#include "clmctfy_mock.h"

struct exec_result {
  int error_code;
  char *message;
  char *function_name;
  struct exec_result *next;
};

extern struct exec_result *result_list;

#define MOCK_FUNCTION_BEGIN(NAME, ...) \
    int NAME (__VA_ARGS__, struct status *s) { \
      if (result_list == NULL || \
          strncmp(result_list->function_name, \
                  #NAME, \
                  strlen(#NAME)) != 0) { \
        lmctfy_mock_set_error_message(#NAME " should not be called. %s should be called instead\n", result_list->function_name); \
        s->message = NULL;  \
        return 0; \
      } \
      struct exec_result *r = result_list;  \
      if (r->error_code == 0) {  \
        s->error_code = 0;  \
        s->message = NULL;  \
      } else {  \
        s->error_code = r->error_code;  \
        s->message = r->message;  \
      } \
      result_list = result_list->next;  \
      free(r->function_name); \
      free(r);  \
      if (s->error_code != 0) {  \
        return s->error_code; \
      }

#define MOCK_FUNCTION_END   \
      return s->error_code; \
    }


#endif // MOCK_CLMCTFY_MACROS_H
