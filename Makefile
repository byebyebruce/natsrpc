types:
	protoc --proto_path=. \
	--proto_path=./third_party \
	--go_out=paths=source_relative:. \
	natsrpc.proto

install:
	go install ./cmd/protoc-gen-natsrpc

test:
	go test ./...
