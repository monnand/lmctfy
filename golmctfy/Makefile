PROTOC = protoc

LMCTFY_PB_GO = $(abspath ./src/containers_lmctfy/lmctfy.pb.go)
LMCTFY_PROTO = $(abspath ../include/lmctfy.proto)

CODES_PB_GO = $(abspath ./src/util_error/codes.pb.go)
CODES_PROTO = $(abspath ../util/task/codes.proto)
PWD=$(abspath .)

# XXX(monnand): Is this an idiomatic way to use protobuf in Go?
build: proto
	GOPATH=${GOPATH}:$(PWD) go build

proto: $(LMCTFY_PB_GO) $(CODES_PB_GO)

$(LMCTFY_PB_GO): $(LMCTFY_PROTO)
	mkdir -p $(dir $(LMCTFY_PB_GO))
	$(PROTOC) --go_out=$(dir $(LMCTFY_PB_GO)) --proto_path=$(dir $(LMCTFY_PROTO)) $(LMCTFY_PROTO)

$(CODES_PB_GO): $(CODES_PROTO)
	mkdir -p $(dir $(CODES_PB_GO))
	$(PROTOC) --go_out=$(dir $(CODES_PB_GO)) --proto_path=$(dir $(CODES_PROTO)) $(CODES_PROTO)

clean:
	rm -rf src/
