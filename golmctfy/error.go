package golmctfy

// #cgo LDFLAGS: -lclmctfy -lprotobuf-c -lprotobuf -lz -lpthread -pthread -lrt -lre2 -lgflags -lstdc++ -lm -L../bin
// #cgo CFLAGS: -I../ -I../include -I../clmctfy/include
// #include "clmctfy.h"
// #include "clmctfy-raw.h"
// #include <stdlib.h>
import "C"
import "unsafe"

type errorStatus struct {
	errorCode    int
	errorMessage string
}

func (self *errorStatus) Error() string {
	return self.errorMessage
}

type Status interface {
	ErrorCode() int
	Error() string
}

func cStatusToGoStatus(s *C.struct_status) error {
	if s != nil && s.error_code != 0 {
		err := new(errorStatus)
		err.errorCode = int(s.error_code)
		if s.message != nil {
			err.errorMessage = C.GoString(s.message)
			C.free(unsafe.Pointer(s.message))
		}
		return err
	}
	return nil
}
