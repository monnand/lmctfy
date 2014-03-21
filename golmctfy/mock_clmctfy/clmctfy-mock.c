#include "clmctfy-mock.h"
#include "macros.h"
#include <string.h>

struct exec_result *result_list = NULL;

void clear_all_expected_calls() {
	struct exec_result *iter = result_list;
	while (iter != NULL) {
		struct exec_result *ptr = iter;
		iter = ptr->next;
		if (ptr->message != NULL) {
			free(ptr->message);
		}
		free(ptr);
	}
	result_list = NULL;
}

void expect_call(void *fn, int error_code, const char *message) {
	struct exec_result *r = (struct exec_result *)malloc(sizeof(struct exec_result));
	r->function_ptr = fn;
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

