# --- Build Stage ---
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth_service ./cmd/auth_service

FROM scratch

COPY --from=builder /app/configs/ /configs/
COPY --from=builder /auth_service /auth_service

EXPOSE 8080

ENTRYPOINT ["/auth_service"]
