package golmctfy

// #cgo pkg-config: protobuf
// #cgo LDFLAGS: -lclmctfy -lprotobuf-c -lprotobuf -lz -lpthread -pthread -lrt -lre2 -lgflags -lstdc++ -lm -L../bin
// #cgo CFLAGS: -I../ -I../include -I../clmctfy/include
// #include <stdlib.h>
// #include <unistd.h>
// #include "clmctfy.h"
// #include "clmctfy-raw.h"
import "C"
import (
	. "containers_lmctfy"
	"fmt"
	"unsafe"
)

const (
	CONTAINER_UPDATE_POLICY_DIFF = iota
	CONTAINER_UPDATE_POLICY_REPLACE
)
const (
	CONTAINER_LIST_POLICY_SELF = iota
	CONTAINER_LIST_POLICY_RECURSIVE
)

const (
	CONTAINER_STATS_TYPE_SUMMARY = iota
	CONTAINER_STATS_TYPE_FULL
)

type Container struct {
	container *C.struct_container
}

// Close the container. Any resource used by the Container object will be
// released. But the underlying container will not be affected.
// A closed container object cannot be used to do anything.
// To destroy the underlying container, one needs to call
// ContainerApi.Destroy()
func (self *Container) Close() error {
	if self == nil || self.container == nil {
		return nil
	}
	C.lmctfy_delete_container(self.container)
	self.container = nil
	return nil
}

func (self *Container) Enter(tids []int) error {
	if self == nil || self.container == nil {
		return ErrInvalidContainer
	}
	if len(tids) == 0 {
		return nil
	}
	ctids := make([]C.pid_t, len(tids))
	for i, tid := range tids {
		ctids[i] = C.pid_t(tid)
	}

	var cstatus C.struct_status
	cstatus.error_code = 0
	C.lmctfy_container_enter(&cstatus, self.container, &ctids[0], C.int(len(ctids)))
	err := cStatusToGoStatus(&cstatus)
	return err
}

func (self *Container) Run(args []string, spec *RunSpec) (tid int, err error) {
	if len(args) == 0 {
		return
	}
	data, size, err := marshalToCData(spec)
	if err != nil {
		return
	}
	cargs := make([]*C.char, len(args))
	for i, arg := range args {
		cargs[i] = C.CString(arg)
	}
	defer func() {
		for _, arg := range cargs {
			C.free(unsafe.Pointer(arg))
		}
	}()

	var ctid C.pid_t
	var cstatus C.struct_status
	cstatus.error_code = 0
	C.lmctfy_container_run_raw(&cstatus, &ctid, self.container, C.int(len(cargs)), &cargs[0], data, size)
	err = cStatusToGoStatus(&cstatus)
	if err != nil {
		return
	}
	tid = int(ctid)
	return
}

func (self *Container) Update(policy int, spec *ContainerSpec) error {
	if self == nil || self.container == nil {
		return ErrInvalidContainer
	}
	data, size, err := marshalToCData(spec)
	if err != nil {
		return err
	}
	var cpolicy C.int
	switch policy {
	case CONTAINER_UPDATE_POLICY_DIFF:
		cpolicy = C.int(C.CONTAINER_UPDATE_POLICY_DIFF)
	case CONTAINER_UPDATE_POLICY_REPLACE:
		cpolicy = C.int(C.CONTAINER_UPDATE_POLICY_REPLACE)
	default:
		return fmt.Errorf("Unknown update policy: %v", policy)
	}
	var cstatus C.struct_status
	cstatus.error_code = 0
	C.lmctfy_container_update_raw(&cstatus, self.container, cpolicy, data, size)
	err = cStatusToGoStatus(&cstatus)
	return err
}
