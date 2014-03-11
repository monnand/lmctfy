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
	"os"
	"reflect"
	"unsafe"

	"github.com/google/lmctfy/golmctfy"
)

func expectCall(fn unsafe.Pointer, code int, msg string) {
	if code == 0 {
		C.expect_call(fn, C.int(0), nil)
		return
	}
	m := C.CString(msg)
	C.expect_call(fn, C.int(code), m)
	C.free(unsafe.Pointer(m))
}

func main() {
	expectCall(C.lmctfy_init_machine_raw, 0, "")
	code := 1
	msg := "error message"
	expectCall(C.lmctfy_init_machine_raw, code, msg)

	var spec InitSpec
	err := golmctfy.InitMachine(&spec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	err = golmctfy.InitMachine(&spec)
	if status, ok := err.(golmctfy.Status); ok {
		if status.ErrorCode() != code {
			fmt.Fprintf(os.Stderr, "Error code should be %v; returned %v\n", code, status.ErrorCode())
		}
	} else {
		fmt.Fprintf(os.Stderr, "Returned type is not a Status, but a %v: %v\n", reflect.TypeOf(err), err)
	}
}
