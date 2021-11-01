gen:
	go install ./tool/cmd/protoc-gen-natsrpc

test:
	go test ./...

serve:
	go run ./tool/cmd/simple_natsserver