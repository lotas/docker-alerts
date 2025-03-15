# Docker Event Listener

Alerts sent to Telegram, Slack or Email on certain docker events.

Can be used as a monitoring tool to know when containers stop working.


## Using in production

```bash
docker run \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e DA_SLACK_WEBHOOK_URL=https://...
  -e DA_TELEGRAM_TOKEN=111:xxxx \
  -e DA_TELEGRAM_CHAT_ID=12345 \
  lotas/docker-alerts
```

## Local development

```bash
make build
make run
make run-debug
make test
```


## Running

1. Using command line flags:
```bash
./docker-alerts
```

2. Using environment variables:
```bash
export DA_SLACK_WEBHOOK_URL=https://...
./docker-alerts
```

3. Using docker (compose)
```bash
docker compose build
docker compose up -d
```


## Credits

Idea borrowed from <https://github.com/luc-ass/docker-telegram-notifier>
with intent to have something written in GO for practicing purposes.
