.PHONY: build

AUTH_BINARY_NAME=auth

auth-build:
	@echo "Building auth application..."
	@go build -o ./build/$(AUTH_BINARY_NAME) ./cmd/main.go

run-auth-local: auth-build
	@echo "Running in LOCAL mode..."
	@./build/$(AUTH_BINARY_NAME) --env-path=.env --config-path=./config/local.yaml
