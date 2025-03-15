package notifications

import (
	"github.com/lotas/docker-alerts/internal/config"
)

func CreateNotifier(cfg *config.Config) Notifier {
	var notifiers []Notifier
	var base []Notifier

	consoleNotifier := NewConsoleNotifier("Docker",
		WithColor(),
	)
	base = append(base, consoleNotifier)

	if cfg.SlackWebhookURL != "" {
		slackNotifier := NewSlackNotifier(
			cfg.SlackWebhookURL,
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
		if cfg.NoDebounce {
			base = append(base, NewMultiNotifier(notifiers...))
		} else {
			// only wrap external api notifiers with debouncer
			// leaving console one as is
			base = append(base, NewDebouncerNotifier(NewMultiNotifier(notifiers...), cfg.DebounceDuration()))
		}
	}

	var notifier Notifier
	if len(base) == 1 {
		notifier = base[0]
	} else {
		notifier = NewMultiNotifier(base...)
	}

	return notifier
}
