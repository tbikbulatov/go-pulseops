.PHONY: infra-up infra-down run-api run-worker run-grpc run-realtime test migrate

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

migrate:
	@echo "migrations are not implemented yet"
