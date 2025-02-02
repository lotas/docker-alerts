package notifications

import (
	"context"

	"github.com/lotas/docker-alerts/internal/config"
)

type Notification struct {
	Title   string
	Message string
	Level   string
}

type Notifier interface {
	Notify(ctx context.Context, notification Notification) error
}

func CreateNotifier(cfg *config.Config) Notifier {
	var notifiers []Notifier

	consoleNotifier := NewConsoleNotifier("DOCKER-ALERT",
		WithColor(),
		WithVerbose(),
	)
	notifiers = append(notifiers, consoleNotifier)

	if cfg.SlackWebhookURL != "" {
		slackNotifier := NewSlackNotifier(
			cfg.SlackWebhookURL,
			cfg.SlackChannel,
		)
		notifiers = append(notifiers, slackNotifier)
	}

	if cfg.TelegramToken != "" && cfg.TelegramChatID != "" {
		telegramNotifier := NewTelegramNotifier(
			cfg.TelegramToken,
			cfg.TelegramChatID,
		)
		notifiers = append(notifiers, telegramNotifier)
	}

	if cfg.EmailSMTPHost != "" {
		emailNotifier := NewEmailNotifier(
			cfg.EmailSMTPHost,
			cfg.EmailSMTPPort,
			cfg.EmailFrom,
			cfg.EmailTo,
		)

		if cfg.EmailSMTPUsername != "" && cfg.EmailSMTPPassword != "" {
			emailNotifier.SetAuth(
				cfg.EmailSMTPUsername,
				cfg.EmailSMTPPassword,
			)
		}

		notifiers = append(notifiers, emailNotifier)
	}

	if len(notifiers) > 1 {
		return NewMultiNotifier(notifiers...)
	}

	return consoleNotifier
}
