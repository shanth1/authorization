.PHONY: build

AUTH_BINARY_NAME=auth
CLIENT_BINARY_NAME=client

# ===== BACKEND =====
run-auth-local: auth-build
	@echo "Running in LOCAL mode..."
	@./build/$(AUTH_BINARY_NAME) --env-path=.env --config-path=./config/local.yaml

run-client-example: client-build
	@./build/$(CLIENT_BINARY_NAME) --env-path=.client.env --config-path=./client/config/config.yaml
# ==================

# ==== FRONTEND =====
run-sso-dev:
	@cd frontend && pnpm run dev:sso

run-portal-dev:
	@cd frontend && pnpm run dev:portal
# ===================

# ===== COMMON =====
auth-build:
	@echo "Building auth application..."
	@go build -o ./build/$(AUTH_BINARY_NAME) ./cmd/main.go

client-build:
	@echo "Building client application..."
	@go build -o ./build/$(CLIENT_BINARY_NAME) ./client/cmd/main.go
# ==================
