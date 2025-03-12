package notifications

import (
	"context"

	"github.com/lotas/docker-alerts/internal/config"
)

type Notifier interface {
	Notify(ctx context.Context, notification Notification, debug bool) error
	NotifyMultiple(ctx context.Context, notifications []Notification, debug bool) error
}

func CreateNotifier(cfg *config.Config) Notifier {
	var notifiers []Notifier
	var base []Notifier

	consoleNotifier := NewConsoleNotifier("DOCKER-ALERT",
		WithColor(),
		WithVerbose(),
	)
	base = append(base, consoleNotifier)

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

	if len(notifiers) > 0 {
		base = append(base, NewMultiNotifier(notifiers...))
	}

	var notifier Notifier
	if len(base) == 1 {
		notifier = base[0]
	} else {
		notifier = NewMultiNotifier(base...)
	}

	if cfg.NoDebounce {
		return notifier
	}

	return NewDebouncerNotifier(notifier, cfg.DebounceDuration())
}
