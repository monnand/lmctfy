package golmctfy

// #cgo LDFLAGS: -lclmctfy -lprotobuf-c -lprotobuf -lz -lpthread -pthread -lrt -lre2 -lgflags -lstdc++ -lm -L../bin
// #cgo CFLAGS: -I../ -I../include -I../clmctfy/include
// #include "clmctfy.h"
// #include "clmctfy-raw.h"
import "C"

type ContainerApi struct {
	containerApi *C.struct_container_api
}

func NewContainerApi() (api *ContainerApi, err error) {
	var cstatus C.struct_status
	api = new(ContainerApi)
	cstatus.error_code = 0
	C.lmctfy_new_container_api(&cstatus, &api.containerApi)
	err = cStatusToGoStatus(&cstatus)
	return
}
