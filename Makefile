.PHONY: infra-up infra-down run-api run-worker run-grpc run-realtime test migrate-up migrate-down migrate-status kafka-ping run-publish-alerts-once run-incident-processor install-tools proto proto-lint observability-up observability-down

ifneq (,$(wildcard .env))
include .env
export
endif
POSTGRES_DSN := host=$(PG_HOST) port=$(PG_PORT) user=$(PG_USER) password=$(PG_PASSWORD) dbname=$(PG_DB) sslmode=$(PG_SSLMODE)

GOOSE_VERSION := v3.27.1
BUF_VERSION := v1.60.0
PROTOC_GEN_GO_VERSION := v1.36.11
PROTOC_GEN_GO_GRPC_VERSION := v1.6.2

infra-up:
	docker compose up -d db redis redpanda

infra-down:
	docker compose down

run-api:
	go run ./cmd/pulseops api

run-publish-alerts-once:
	go run ./cmd/pulseops publish-alerts-once

run-incident-processor:
	go run ./cmd/pulseops incident-processor

run-worker:
	go run ./cmd/pulseops

run-grpc:
	go run ./cmd/pulseops grpc

run-realtime:
	go run ./cmd/pulseops

test:
	go test ./...

kafka-ping:
	go run ./cmd/pulseops kafka-ping

install-tools:
	go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)
	go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

migrate-up:
	goose -dir migrations postgres "$(POSTGRES_DSN)" up

migrate-down:
	goose -dir migrations postgres "$(POSTGRES_DSN)" down

migrate-status:
	goose -dir migrations postgres "$(POSTGRES_DSN)" status

proto:
	buf generate

proto-lint:
	buf lint

observability-up:
	docker compose --profile observability up -d prometheus

observability-down:
	docker compose stop prometheus
