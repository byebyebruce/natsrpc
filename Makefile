protoc-gen-natsrpc:
	go install ./tool/cmd/protoc-gen-natsrpc

test:
	go test ./...

serve:
	go run ./tool/cmd/simple_natsserver

example_pb:
	sh example/pb/gen.sh
	sh example/pb/request/gen.sh
	sh example/pb/publish/gen.sh
	sh example/pb/async_service/gen.sh
	sh example/pb/async_client/gen.sh
