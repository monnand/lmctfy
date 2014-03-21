package golmctfy

// #cgo pkg-config: protobuf
// #cgo LDFLAGS: -lclmctfy -lprotobuf-c -lprotobuf -lz -lpthread -pthread -lrt -lre2 -lgflags -lstdc++ -lm -L../../bin
// #cgo CFLAGS: -I../ -I../include
// #include <stdlib.h>
// #include <unistd.h>
// #include "clmctfy.h"
// #include "clmctfy-raw.h"
//
// extern void lmctfy_mock_assert_expectations();
// extern void lmctfy_mock_expect_call(const char *fn, int error_code, const char *message);
// extern const char *lmctfy_mock_get_last_error_message();
// extern void lmctfy_mock_clear_last_error_message();
import "C"

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

func expectCall(fn string, code int, msg string, action func() error) error {
	fn_cstr := C.CString(fn)
	defer C.free(unsafe.Pointer(fn_cstr))
	if code == 0 {
		C.lmctfy_mock_expect_call(fn_cstr, C.int(0), nil)
	} else {
		m := C.CString(msg)
		C.lmctfy_mock_expect_call(fn_cstr, C.int(code), m)
		C.free(unsafe.Pointer(m))
	}
	err := action()
	last_errmsg := C.lmctfy_mock_get_last_error_message()
	if last_errmsg != nil {
		defer C.lmctfy_mock_clear_last_error_message()
		return errors.New(C.GoString(last_errmsg))
	}
	if code == 0 {
		if err != nil {
			return fmt.Errorf("Should return successfully. But received error: %v\n", err)
		}
		return nil
	}
	if status, ok := err.(Status); ok {
		if status.ErrorCode() != code {
			return fmt.Errorf("Error code should be %v; returned %v\n", code, status.ErrorCode())
		}
		if status.Error() != msg {
			return fmt.Errorf("Error message should be %v; returned %v\n", msg, status.Error())
		}
	} else {
		return fmt.Errorf("Returned type is not a Status, but a %v: %v\n", reflect.TypeOf(err), err)
	}
	return nil
}

func assertExpectations() error {
	C.lmctfy_mock_assert_expectations()
	last_errmsg := C.lmctfy_mock_get_last_error_message()
	if last_errmsg == nil {
		return nil
	}
	defer C.lmctfy_mock_clear_last_error_message()

	return errors.New(C.GoString(last_errmsg))
}
