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

func (s *SlackNotifier) Notify(ctx context.Context, notification Notification, debug bool) error {
	msg := slack.WebhookMessage{
		Channel: s.channel,
		Text:    notification.Message,
	}

	return slack.PostWebhook(s.webhookURL, &msg)
}

func (c *SlackNotifier) NotifyMultiple(ctx context.Context, notifications []Notification, debug bool) error {
	for _, n := range notifications {
		c.Notify(ctx, n, debug)
	}
	return nil
}
