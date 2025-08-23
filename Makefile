.PHONY: build

AUTH_BINARY_NAME=auth
CLIENT_BINARY_NAME=client

auth-build:
	@echo "Building auth application..."
	@go build -o ./build/$(AUTH_BINARY_NAME) ./cmd/main.go

run-auth-local: auth-build
	@echo "Running in LOCAL mode..."
	@./build/$(AUTH_BINARY_NAME) --env-path=.env --config-path=./config/local.yaml

client-build:
	@echo "Building client application..."
	@go build -o ./build/$(CLIENT_BINARY_NAME) ./client/cmd/main.go

run-client-example: client-build
	@./build/$(CLIENT_BINARY_NAME) --env-path=.client.env --config-path=./client/config/config.yaml
