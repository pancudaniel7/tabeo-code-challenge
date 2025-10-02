# Makefile for tabeo-code-challenge

.PHONY: help build run test docker-up docker-down

help:
	@echo "Available targets:"
	@echo "  help         Show this help message"
	@echo "  build        Build the Go application"
	@echo "  run          Run the Go application"
	@echo "  test         Run all tests"
	@echo "  docker-up    Start docker-compose services (e.g., DB, Redis)"
	@echo "  docker-down  Stop docker-compose services"

build:
	go build -o bin/tabeo ./cmd/tabeo

run: build
	./bin/tabeo

test:
	go test ./...

docker-up:
	docker-compose -f deployments/docker-compose.yml up -d

docker-down:
	docker-compose -f deployments/docker-compose.yml down

