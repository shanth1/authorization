.PHONY: build

AUTH_BINARY_NAME=auth
BOT_BINARY_NAME=bot

auth-build:
	@echo "Building auth application..."
	@go build -o ./build/$(AUTH_BINARY_NAME) ./cmd/auth/main.go

bot-build:
	@echo "Building bot application..."
	@go build -o ./build/$(BOT_BINARY_NAME) ./cmd/bot/main.go

run-auth-local: auth-build
	@echo "Running in LOCAL mode..."
	@./build/$(AUTH_BINARY_NAME) --env-path=.env --config-path=./config/auth/local.yaml

run-bot-local: bot-build
	@echo "Running BOT in LOCAL mode..."
	@./build/$(BOT_BINARY_NAME) --env-path=.env
