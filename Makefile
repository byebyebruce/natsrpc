all: types protoc-gen-natsrpc  example test

types:
	protoc --proto_path=.:./third_party \
	--go_out=paths=source_relative:. \
	natsrpc.proto

protoc-gen-natsrpc:
	go install ./cmd/protoc-gen-natsrpc

test:
	go test ./...

.PHONY: example
example:
	make -C example
