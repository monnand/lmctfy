#ifndef MOCK_CLMCTFY_MOCK_H_
#define MOCK_CLMCTFY_MOCK_H_

void lmctfy_mock_expect_call(const char *fn, int error_code, const char *message);
void lmctfy_mock_clear_all_expected_calls();
void lmctfy_mock_set_error_message(const char *fmt, ...);

#endif // MOCK_CLMCTFY_MOCK_H_
