#ifndef MOCK_CLMCTFY_MOCK_H_
#define MOCK_CLMCTFY_MOCK_H_

#include "clmctfy.h"

#define CONTAINER_API_ADDR	(struct container_api *)0x1234

struct container_api {
};

struct container {
	char *name;
	void *cb_userdata;
	lmctfy_event_callback_f callback;
};

void lmctfy_mock_expect_call(const char *fn, int error_code, const char *message);
void lmctfy_mock_clear_all_expected_calls();
void lmctfy_mock_set_error_message(const char *fmt, ...);

#endif // MOCK_CLMCTFY_MOCK_H_
