all: generate protoc-gen-natsrpc  example test

generate:
	go generate ./...

protoc-gen-natsrpc:
	go install ./cmd/protoc-gen-natsrpc

test:
	go test ./...

.PHONY: example
example:
	make -C example
