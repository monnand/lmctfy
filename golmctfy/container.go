package golmctfy

// #cgo pkg-config: protobuf
// #cgo LDFLAGS: -lclmctfy -lprotobuf-c -lprotobuf -lz -lpthread -pthread -lrt -lre2 -lgflags -lstdc++ -lm -L../bin
// #cgo CFLAGS: -I../ -I../include
// #include <stdlib.h>
// #include <unistd.h>
// #include "clmctfy.h"
// #include "clmctfy-raw.h"
// extern void golmctfy_cgo_notif_callback(struct container *container,
//                                  const struct status *status,
//                                  void *userdata);
import "C"
import (
	. "containers_lmctfy"
	"errors"
	"fmt"
	"reflect"
	"sync"
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

type Event struct {
	Container *Container
	NotifId   uint64
	Error     error
}

type userData struct {
	container *Container
	notifId   uint64
	ch        chan<- *Event
}

type Container struct {
	container *C.struct_container
	lock      sync.RWMutex
	// store user data so that it will not be released by GC.
	userDataMap map[uint64]*userData
}

//export golmctfyNotifCallback
func golmctfyNotifCallback(status *C.struct_status, ptr unsafe.Pointer) {
	var udata *userData
	udata = (*userData)(ptr)
	var err error
	err = nil
	if status.error_code != 0 {
		e := new(errorStatus)
		e.errorCode = int(status.error_code)
		if status.message != nil {
			e.errorMessage = C.GoString(status.message)
		}
		err = e
	}
	evt := &Event{
		Container: udata.container,
		NotifId:   udata.notifId,
		Error:     err,
	}
	go func() {
		udata.ch <- evt
	}()
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
	C.lmctfy_container_enter(self.container, &ctids[0], C.int(len(ctids)), &cstatus)
	err := cStatusToGoStatus(&cstatus)
	return err
}

func (self *Container) Run(args []string, spec *RunSpec) (tid int, err error) {
	if self == nil || self.container == nil {
		err = ErrInvalidContainer
		return
	}
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
	C.lmctfy_container_run_raw(self.container, C.int(len(cargs)), &cargs[0], data, size, &ctid, &cstatus)
	err = cStatusToGoStatus(&cstatus)
	if err != nil {
		return
	}
	tid = int(ctid)
	return
}

func (self *Container) Exec(args []string) error {
	if self == nil || self.container == nil {
		return ErrInvalidContainer
	}
	if len(args) == 0 {
		return errors.New("Not enough arguments")
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

	var cstatus C.struct_status
	C.lmctfy_container_exec(self.container, C.int(len(cargs)), &cargs[0], &cstatus)
	return cStatusToGoStatus(&cstatus)
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
	C.lmctfy_container_update_raw(self.container, cpolicy, data, size, &cstatus)
	err = cStatusToGoStatus(&cstatus)
	return err
}

func (self *Container) Name() string {
	if self == nil || self.container == nil {
		return ""
	}
	cname := C.lmctfy_container_name(self.container)
	if cname == nil {
		return ""
	}
	return C.GoString(cname)
}

// ch: ch will nevery be closed by the container even if the notification is
// unregistered.
func (self *Container) RegisterNotification(spec *EventSpec, ch chan<- *Event) (notifId uint64, err error) {
	if self == nil || self.container == nil {
		err = ErrInvalidContainer
		return
	}
	if ch == nil {
		err = errors.New("Invalid channel")
		return
	}
	data, size, err := marshalToCData(spec)
	if err != nil {
		return
	}
	var cstatus C.struct_status
	cstatus.error_code = 0

	ud := new(userData)
	ud.ch = ch
	ud.container = self
	var nid C.notification_id_t
	self.lock.Lock()
	defer self.lock.Unlock()
	C.lmctfy_container_register_notification_raw(self.container,
		(*[0]byte)(C.golmctfy_cgo_notif_callback),
		unsafe.Pointer(ud),
		data,
		size,
		&nid,
		&cstatus)
	err = cStatusToGoStatus(&cstatus)
	if err != nil {
		return
	}
	notifId = uint64(nid)
	ud.notifId = notifId
	if self.userDataMap == nil {
		self.userDataMap = make(map[uint64]*userData, 10)
	}
	self.userDataMap[notifId] = ud
	return
}

// The corresponding channel will not be closed by this method. It's the
// caller's responsibility to properly close the channel.
func (self *Container) UnregisterNotification(notifId uint64) error {
	if self == nil || self.container == nil {
		return ErrInvalidContainer
	}
	var cstatus C.struct_status
	cstatus.error_code = 0

	self.lock.Lock()
	defer self.lock.Unlock()
	if self.userDataMap != nil {
		if _, ok := self.userDataMap[notifId]; ok {
			delete(self.userDataMap, notifId)
		}
	}
	C.lmctfy_container_unregister_notification(self.container,
		C.notification_id_t(notifId),
		&cstatus)
	return cStatusToGoStatus(&cstatus)
}

func (self *Container) Pause() error {
	if self == nil || self.container == nil {
		return ErrInvalidContainer
	}
	var cstatus C.struct_status
	cstatus.error_code = 0

	C.lmctfy_container_pause(self.container, &cstatus)
	return cStatusToGoStatus(&cstatus)
}

func (self *Container) Resume() error {
	if self == nil || self.container == nil {
		return ErrInvalidContainer
	}
	var cstatus C.struct_status
	cstatus.error_code = 0

	C.lmctfy_container_resume(self.container, &cstatus)
	return cStatusToGoStatus(&cstatus)
}

func (self *Container) KillAll() error {
	if self == nil || self.container == nil {
		return ErrInvalidContainer
	}
	var cstatus C.struct_status
	cstatus.error_code = 0

	C.lmctfy_container_killall(self.container, &cstatus)
	return cStatusToGoStatus(&cstatus)
}

func (self *Container) ListThreads(policy int) (threads []int, err error) {
	if self == nil || self.container == nil {
		err = ErrInvalidContainer
		return
	}
	var cstatus C.struct_status
	cstatus.error_code = 0

	var pids *C.pid_t
	var n C.int
	var p C.int

	switch policy {
	case CONTAINER_LIST_POLICY_SELF:
		p = C.CONTAINER_LIST_POLICY_SELF
	case CONTAINER_LIST_POLICY_RECURSIVE:
		p = C.CONTAINER_LIST_POLICY_RECURSIVE
	default:
		err = fmt.Errorf("Unknown list policy: %v", policy)
		return
	}

	C.lmctfy_container_list_threads(self.container, p, &pids, &n, &cstatus)
	err = cStatusToGoStatus(&cstatus)
	if err != nil {
		return
	}
	defer C.free(unsafe.Pointer(pids))
	nrPids := int(n)
	threads = make([]int, nrPids)
	var cthreads []C.pid_t
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&cthreads)))
	sliceHeader.Cap = nrPids
	sliceHeader.Len = nrPids
	sliceHeader.Data = uintptr(unsafe.Pointer(pids))

	for i, pid := range cthreads {
		threads[i] = int(pid)
	}
	return
}

func (self *Container) ListProcesses(policy int) (processes []int, err error) {
	if self == nil || self.container == nil {
		err = ErrInvalidContainer
		return
	}
	var cstatus C.struct_status
	cstatus.error_code = 0

	var pids *C.pid_t
	var n C.int
	var p C.int

	switch policy {
	case CONTAINER_LIST_POLICY_SELF:
		p = C.CONTAINER_LIST_POLICY_SELF
	case CONTAINER_LIST_POLICY_RECURSIVE:
		p = C.CONTAINER_LIST_POLICY_RECURSIVE
	default:
		err = fmt.Errorf("Unknown list policy: %v", policy)
		return
	}

	C.lmctfy_container_list_processes(self.container, p, &pids, &n, &cstatus)
	err = cStatusToGoStatus(&cstatus)
	if err != nil {
		return
	}
	defer C.free(unsafe.Pointer(pids))
	nrPids := int(n)
	processes = make([]int, nrPids)
	var cpids []C.pid_t
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&cpids)))
	sliceHeader.Cap = nrPids
	sliceHeader.Len = nrPids
	sliceHeader.Data = uintptr(unsafe.Pointer(pids))

	for i, pid := range cpids {
		processes[i] = int(pid)
	}
	return
}
