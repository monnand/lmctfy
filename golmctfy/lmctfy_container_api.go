package golmctfy

// #cgo LDFLAGS: -lclmctfy -lprotobuf-c -lprotobuf -lz -lpthread -pthread -lrt -lre2 -lgflags -lstdc++ -lm -L../bin
// #cgo CFLAGS: -I../ -I../include -I../clmctfy/include
// #include "clmctfy.h"
// #include "clmctfy-raw.h"
import "C"
import "errors"

type ContainerApi struct {
	containerApi *C.struct_container_api
}

type errorStatus struct {
	errorCode    int
	errorMessage string
}

type Status interface {
	ErrorCode() int
	ErrorMessage() string
}

func NewContainerApi() (api *ContainerApi, err error) {
	var cstatus C.struct_status
	api = new(ContainerApi)
	cstatus.error_code = 0
	C.lmctfy_new_container_api(&cstatus, &api.containerApi)
	if cstatus.error_code != 0 {
		api = nil
		err = errors.New(C.GoString(cstatus.message))
		return
	}
	err = nil
	return
}
