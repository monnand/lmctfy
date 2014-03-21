package golmctfy

// #cgo pkg-config: protobuf
// #cgo LDFLAGS: -lclmctfy -lprotobuf-c -lprotobuf -lz -lpthread -pthread -lrt -lre2 -lgflags -lstdc++ -lm -L../../bin
// #cgo CFLAGS: -I../ -I../include
// #include <stdlib.h>
// #include <unistd.h>
// #include "clmctfy.h"
// #include "clmctfy-raw.h"
//
// extern void lmctfy_mock_expect_call(const char *fn, int error_code, const char *message);
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

func expectCall(fn string, code int, msg string, action func() error) {
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
	if code == 0 {
		if err != nil {
			panic(fmt.Sprintf("Should return successfully. But received error: %v\n", err))
			return
		}
		return
	}
	if status, ok := err.(Status); ok {
		if status.ErrorCode() != code {
			panic(fmt.Sprintf("Error code should be %v; returned %v\n", code, status.ErrorCode()))
		}
		if status.Error() != msg {
			panic(fmt.Sprintf("Error message should be %v; returned %v\n", msg, status.Error()))
		}
	} else {
		panic(fmt.Sprintf("Returned type is not a Status, but a %v: %v\n", reflect.TypeOf(err), err))
	}
}
