package notifications

import (
	"context"
	"github.com/slack-go/slack"
)

type SlackNotifier struct {
	webhookURL string
}

func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: webhookURL,
	}
}

func (s *SlackNotifier) Notify(ctx context.Context, event Event, debug bool) error {
	msg := slack.WebhookMessage{
		Text: event.Markdown(),
	}

	return slack.PostWebhook(s.webhookURL, &msg)
}

func (c *SlackNotifier) NotifyMultiple(ctx context.Context, events []Event, debug bool) error {
	for _, n := range events {
		if err := c.Notify(ctx, n, debug); err != nil {
			return err
		}
	}
	return nil
}
