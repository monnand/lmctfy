package golmctfy

// #cgo pkg-config: protobuf
// #cgo LDFLAGS: -lclmctfy -lprotobuf-c -lprotobuf -lz -lpthread -pthread -lrt -lre2 -lgflags -lstdc++ -lm -L../bin
// #cgo CFLAGS: -I../ -I../include
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
	if len(d) == 0 {
		data = nil
		size = C.size_t(0)
		err = nil
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
	C.lmctfy_init_machine_raw(data, size, &cstatus)
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
	C.lmctfy_new_container_api(&api.containerApi, &cstatus)
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
	var c *C.struct_container
	str := C.CString(container_name)
	defer C.free(unsafe.Pointer(str))
	C.lmctfy_container_api_create_container_raw(self.containerApi, str, data, size, &c, &cstatus)
	err = cStatusToGoStatus(&cstatus)
	if err != nil {
		container = nil
		return
	}
	container, err = newContainer(c)
	return
}

func (self *ContainerApi) Get(container_name string) (container *Container, err error) {
	if self == nil || self.containerApi == nil {
		err = ErrInvalidContainerApi
		return
	}
	var cstatus C.struct_status
	cstatus.error_code = 0
	var c *C.struct_container
	str := C.CString(container_name)
	defer C.free(unsafe.Pointer(str))

	C.lmctfy_container_api_get_container(self.containerApi, str, &c, &cstatus)
	err = cStatusToGoStatus(&cstatus)
	if err != nil {
		container = nil
		return
	}
	container, err = newContainer(c)
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

	C.lmctfy_container_api_destroy_container(self.containerApi, container.container, &cstatus)
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

	C.lmctfy_container_api_detect_container(self.containerApi, cpid, &cname, &cstatus)

	err = cStatusToGoStatus(&cstatus)
	if err != nil {
		return
	}
	container_name = C.GoString(cname)
	C.free(unsafe.Pointer(cname))
	return
}
