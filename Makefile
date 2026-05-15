.PHONY: infra-up infra-down run-api run-worker run-grpc run-realtime test migrate-up migrate-down migrate-status

ifneq (,$(wildcard .env))
include .env
export
endif
POSTGRES_DSN := host=$(PG_HOST) port=$(PG_PORT) user=$(PG_USER) password=$(PG_PASSWORD) dbname=$(PG_DB) sslmode=$(PG_SSLMODE)

infra-up:
	docker compose up -d db redis redpanda

infra-down:
	docker compose down

run-api:
	go run ./cmd/pulseops

run-worker:
	go run ./cmd/pulseops

run-grpc:
	go run ./cmd/pulseops

run-realtime:
	go run ./cmd/pulseops

test:
	go test ./...

install-tools:
	go install github.com/pressly/goose/v3/cmd/goose@latest

migrate-up:
	goose -dir migrations postgres "$(POSTGRES_DSN)" up

migrate-down:
	goose -dir migrations postgres "$(POSTGRES_DSN)" down

migrate-status:
	goose -dir migrations postgres "$(POSTGRES_DSN)" status
