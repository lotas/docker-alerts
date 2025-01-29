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

func (s *SlackNotifier) Notify(ctx context.Context, notification Notification) error {
	msg := slack.WebhookMessage{
		Channel: s.channel,
		Text:    notification.Message,
	}

	return slack.PostWebhook(s.webhookURL, &msg)
}
