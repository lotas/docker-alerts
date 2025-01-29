# Variables
BINARY_NAME=docker-event-listener
MAIN_PACKAGE=cmd/main.go
DOCKER_COMPOSE=docker-compose.yml

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## build: Build the binary
build:
	@echo "Building..."
	go build -o $(GOBIN)/$(BINARY_NAME) $(MAIN_PACKAGE)

## run: Run the application
run:
	go run $(MAIN_PACKAGE)

## clean: Clean build files
clean:
	@echo "Cleaning..."
	@rm -rf $(GOBIN)
	@go clean

## test: Run tests
test:
	go test ./... -v

## coverage: Run tests with coverage
coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

## lint: Run linter
lint:
	@which golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
	golangci-lint run

## fmt: Format code
fmt:
	go fmt ./...

## vet: Run go vet
vet:
	go vet ./...

## deps: Download dependencies
deps:
	go mod download

## tidy: Tidy go.mod
tidy:
	go mod tidy

## help: Display this help screen
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build run clean test coverage lint fmt vet deps tidy help
