#include "clmctfy-mock.h"
#include "macros.h"
#include <string.h>

struct exec_result *result_list = NULL;

void expect_call(void *fn, int error_code, const char *message) {
	struct exec_result *r = (struct exec_result *)malloc(sizeof(struct exec_result));
	r->function_ptr = fn;
	r->error_code = error_code;
	if (error_code != 0) {
		r->message = strdup(message);
	}
	struct exec_result *iter = NULL;
	if (result_list == NULL) {
		result_list = r;
		return;
	}
	for (iter = result_list; iter->next != NULL; iter = iter->next) {
	}

	iter->next = r;
}

