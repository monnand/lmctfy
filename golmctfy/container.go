package golmctfy

// #cgo pkg-config: protobuf
// #cgo LDFLAGS: -lclmctfy -lprotobuf-c -lprotobuf -lz -lpthread -pthread -lrt -lre2 -lgflags -lstdc++ -lm -L../bin
// #cgo CFLAGS: -I../ -I../include -I../clmctfy/include
// #include <stdlib.h>
// #include "clmctfy.h"
// #include "clmctfy-raw.h"
import "C"

type Container struct {
	container *C.struct_container
}

func (self *Container) Close() error {
	if self == nil || self.container == nil {
		return nil
	}
	C.lmctfy_delete_container(self.container)
	self.container = nil
	return nil
}
