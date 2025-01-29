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
docker-alerts --docker-host="tcp://127.0.0.1:2375" --slack-webhook-url="https://..."
```

2. Using environment variables:
```bash
export DA_DOCKER_HOST=tcp://127.0.0.1:2375
export DA_SLACK_WEBHOOK_URL=https://...
docker-alerts
```

3. Using a config file:
```bash
docker-alerts --config=/path/to/config.yaml
```
