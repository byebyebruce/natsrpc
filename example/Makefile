pb:
	protoc \
    --proto_path=. \
    --proto_path=.. \
    --go_out=paths=source_relative:. \
    --natsrpc_out=paths=source_relative:. \
    *.proto

run:
	go run ./...