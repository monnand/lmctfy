package main

// #cgo pkg-config: protobuf
// #cgo LDFLAGS: -lclmctfy -lprotobuf-c -lprotobuf -lz -lpthread -pthread -lrt -lre2 -lgflags -lstdc++ -lm -L../../bin
// #cgo CFLAGS: -I../../ -I../../include -I../../clmctfy/include -I../mock_clmctfy
// #include <stdlib.h>
// #include <unistd.h>
// #include "clmctfy.h"
// #include "clmctfy-raw.h"
// #include "clmctfy-mock.h"
import "C"

import (
	. "containers_lmctfy"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/google/lmctfy/golmctfy"
)

func expectCall(fn unsafe.Pointer, code int, msg string, action func() error) {
	if code == 0 {
		C.expect_call(fn, C.int(0), nil)
	} else {
		m := C.CString(msg)
		C.expect_call(fn, C.int(code), m)
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
	if status, ok := err.(golmctfy.Status); ok {
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

func main() {
	expectCall(C.lmctfy_new_container_api, 0, "", func() error {
		api, err := golmctfy.NewContainerApi()
		defer api.Close()
		return err
	})
	expectCall(C.lmctfy_new_container_api, 5, "error message", func() error {
		api, err := golmctfy.NewContainerApi()
		defer api.Close()
		return err
	})
}
