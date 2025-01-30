# Docker Event Listener

A Go application that listens to Docker events and triggers actions based on specific events.

## Prerequisites

- Go 1.16 or later
- Docker
- Make

## Installation

```bash
make deps
```

## Usage

```bash
# Build the binary
make build

# Run the application
make run

# Run tests
make test

# See all available commands
make help

## Running

1. Using command line flags:
```bash
docker-alerts
```

2. Using environment variables:
```bash
export DA_SLACK_WEBHOOK_URL=https://...
docker-alerts
```
