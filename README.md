# Docker Event Listener

A Go application that listens to Docker events and triggers actions based on specific events.

## TODO

- [ ] Notifications to
  - [x] Telegram
  - [ ] Slack
  - [ ] Email
- [ ] ShouldNotify - filter which events to send
- [ ] Use docker labels
  - [ ] supress container notifications
  - [ ] use custom notification transport or params
  - [ ] use container filters for events to decide in ShouldNotify


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

3. Using docker (compose)
```bash
docker compose build
docker compose up
```


## Credits

Idea borrowed from <https://github.com/luc-ass/docker-telegram-notifier>
with intent to have something written in GO for practicing purposes.
