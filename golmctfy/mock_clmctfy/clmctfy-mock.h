#ifndef MOCK_CLMCTFY_MOCK_H_
#define MOCK_CLMCTFY_MOCK_H_

void expect_call(void *fn, int error_code, const char *message);
void clear_all_expected_calls();

#endif // MOCK_CLMCTFY_MOCK_H_
