EXAMPLES := 1_request 2_publish 3_header 4_namespace 5_id \
	6_handle_error 7_recovery 8_middleware 9_loadbalance 10_reply_later \
	11_encoder 12_single_goroutine

run:
	@for EXAMPLE in $(EXAMPLES); do \
  		echo "run example: $$EXAMPLE"; \
    	go run ./$$EXAMPLE; \
    done

all:
	make -C ../ install
	make pb run

pb:
	protoc \
    --proto_path=. \
    --proto_path=.. \
	--proto_path=../third_party \
    --gogo_out=paths=source_relative:. \
    --natsrpc_out=paths=source_relative:. \
    *.proto
