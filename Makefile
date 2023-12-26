include .envrc

PROTO_DIR=./api

## help: print this help message
.PHONY: help
help:
	@echo 'usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'are you sure? [y/N]' && read ans && [ $${ans:-n} = y ]

## proto: generate code from proto files
.PHONY: proto
proto:
	@echo "generating code from proto files..."
	@for file in $(PROTO_DIR)/*.proto; do \
		service_name=$$(basename -s .proto $$file); \
		protoc --proto_path=$(PROTO_DIR) \
		--go_out=paths=source_relative:./pkg/api/$$service_name \
		--go-grpc_out=paths=source_relative:./pkg/api/$$service_name \
		$$file; \
	done

## proto/clean: delete code generated from proto files
.PHONY: proto/clean
proto/clean: confirm
	@echo "removing code generated from proto files..."
	find . -type f -name "*.pb.go" | xargs -I {} rm -f {}


## run/chat: run chat service
.PHONY: run/chat
run/chat:
	@echo "Running chat service..."
	go run ./cmd/chat/main.go -address=$(CHAT_SERVICE_ADDRESS)
