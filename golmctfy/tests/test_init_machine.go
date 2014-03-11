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

	"github.com/google/lmctfy/golmctfy"
)

func main() {
	var spec InitSpec
	C.expect_call(C.lmctfy_init_machine_raw, 0, nil)
	err := golmctfy.InitMachine(&spec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
