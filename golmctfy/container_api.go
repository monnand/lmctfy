package golmctfy

// #cgo pkg-config: protobuf
// #cgo LDFLAGS: -lclmctfy -lprotobuf-c -lprotobuf -lz -lpthread -pthread -lrt -lre2 -lgflags -lstdc++ -lm -L../bin
// #cgo CFLAGS: -I../ -I../include -I../clmctfy/include
// #include <unistd.h>
// #include <stdlib.h>
// #include "clmctfy.h"
// #include "clmctfy-raw.h"
import "C"
import (
	. "containers_lmctfy"
	"unsafe"

	"code.google.com/p/goprotobuf/proto"
)

func marshalToCData(spec proto.Message) (data unsafe.Pointer, size C.size_t, err error) {
	d, err := proto.Marshal(spec)
	if err != nil {
		return
	}
	data = unsafe.Pointer(&d[0])
	size = C.size_t(len(d))
	err = nil
	return
}

func InitMachine(spec *InitSpec) error {
	var cstatus C.struct_status
	cstatus.error_code = 0
	data, size, err := marshalToCData(spec)
	if err != nil {
		return err
	}
	C.lmctfy_init_machine_raw(&cstatus, data, size)
	err = cStatusToGoStatus(&cstatus)
	return err
}

type ContainerApi struct {
	containerApi *C.struct_container_api
}

func NewContainerApi() (api *ContainerApi, err error) {
	var cstatus C.struct_status
	cstatus.error_code = 0
	api = new(ContainerApi)
	C.lmctfy_new_container_api(&cstatus, &api.containerApi)
	err = cStatusToGoStatus(&cstatus)
	return
}

func (self *ContainerApi) Close() error {
	if self == nil || self.containerApi == nil {
		return nil
	}

	C.lmctfy_delete_container_api(self.containerApi)
	self.containerApi = nil
	return nil
}

func (self *ContainerApi) Create(container_name string, spec *ContainerSpec) (container *Container, err error) {
	if self == nil || self.containerApi == nil {
		err = ErrInvalidContainerApi
		return
	}
	var cstatus C.struct_status
	cstatus.error_code = 0
	data, size, err := marshalToCData(spec)
	if err != nil {
		return
	}
	container = new(Container)
	str := C.CString(container_name)
	defer C.free(unsafe.Pointer(str))
	C.lmctfy_container_api_create_container_raw(&cstatus, &(container.container), self.containerApi, str, data, size)
	err = cStatusToGoStatus(&cstatus)
	if err != nil {
		container = nil
	}
	return
}

func (self *ContainerApi) Get(container_name string) (container *Container, err error) {
	if self == nil || self.containerApi == nil {
		err = ErrInvalidContainerApi
		return
	}
	var cstatus C.struct_status
	cstatus.error_code = 0
	container = new(Container)
	str := C.CString(container_name)
	defer C.free(unsafe.Pointer(str))

	C.lmctfy_container_api_get_container(&cstatus, &(container.container), self.containerApi, str)
	err = cStatusToGoStatus(&cstatus)
	if err != nil {
		container = nil
	}
	return
}

// Destroy a container. The container object will also be closed
func (self *ContainerApi) Destroy(container *Container) error {
	if self == nil || self.containerApi == nil {
		err := ErrInvalidContainerApi
		return err
	}
	if container == nil || container.container == nil {
		err := ErrInvalidContainer
		return err
	}
	var cstatus C.struct_status
	cstatus.error_code = 0

	C.lmctfy_container_api_destroy_container(&cstatus, self.containerApi, container.container)
	err := cStatusToGoStatus(&cstatus)
	if err == nil {
		container.Close()
	}
	return err
}

func (self *ContainerApi) Detect(pid int) (container_name string, err error) {
	if self == nil || self.containerApi == nil {
		err = ErrInvalidContainerApi
		return
	}
	cpid := C.pid_t(pid)
	var cstatus C.struct_status
	cstatus.error_code = 0

	var cname *C.char

	C.lmctfy_container_api_detect_container(&cstatus, &cname, self.containerApi, cpid)

	err = cStatusToGoStatus(&cstatus)
	if err != nil {
		return
	}
	container_name = C.GoString(cname)
	C.free(unsafe.Pointer(cname))
	return
}
