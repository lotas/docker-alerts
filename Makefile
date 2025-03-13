BINARY_NAME=docker-event-listener
MAIN_PACKAGE=main.go
DOCKER_COMPOSE=docker-compose.yml

GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

MAKEFLAGS += --silent

build:
	@echo "Building..."
	go build -o $(GOBIN)/$(BINARY_NAME) $(MAIN_PACKAGE)

run:
	go run $(MAIN_PACKAGE)

run-debug:
	go run $(MAIN_PACKAGE) --debug --debounce-seconds=10

clean:
	@echo "Cleaning..."
	@rm -rf $(GOBIN)
	@go clean

test:
	go test ./... -v

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

lint:
	@which golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
	golangci-lint run

fmt:
	go fmt ./...

vet:
	go vet ./...

deps:
	go mod download

tidy:
	go mod tidy

.PHONY: build run clean test coverage lint fmt vet deps tidy
