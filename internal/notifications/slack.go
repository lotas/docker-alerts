package notifications

import (
	"context"
	"github.com/slack-go/slack"
)

type SlackNotifier struct {
	client     *slack.Client
	channel    string
	webhookURL string
}

func NewSlackNotifier(webhookURL, channel string) *SlackNotifier {
	return &SlackNotifier{
		client:  slack.New(webhookURL),
		channel: channel,
	}
}

func (s *SlackNotifier) Notify(ctx context.Context, event Event, debug bool) error {
	msg := slack.WebhookMessage{
		Channel: s.channel,
		Text:    event.Text(),
	}

	return slack.PostWebhook(s.webhookURL, &msg)
}

func (c *SlackNotifier) NotifyMultiple(ctx context.Context, events []Event, debug bool) error {
	for _, n := range events {
		c.Notify(ctx, n, debug)
	}
	return nil
}
