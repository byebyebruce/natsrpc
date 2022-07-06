all: generate protoc-gen-natsrpc test example

generate:
	go generate ./...

protoc-gen-natsrpc:
	cd cmd/protoc-gen-natsrpc && go install ./

test:
	go test ./...

.PHONY: example
example:
	make -C example
