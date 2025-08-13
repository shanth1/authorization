.PHONY: build

BINARY_NAME=auth

build:
	@echo "Building auth application..."
	@go build -o ./build/$(BINARY_NAME) ./cmd/auth/main.go

run-local: build
	@echo "Running in LOCAL mode..."
	@./build/$(BINARY_NAME) --env-path=.env --config-path=./config/auth/local.yaml
