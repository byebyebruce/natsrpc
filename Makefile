all: protoc-gen-natsrpc test example

protoc-gen-natsrpc:
	go install ./tool/cmd/protoc-gen-natsrpc

test:
	go test ./...

.PHONY: example
example:
	make -C example
