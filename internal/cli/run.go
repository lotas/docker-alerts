package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lotas/docker-alerts/internal/config"
	"github.com/lotas/docker-alerts/internal/docker"
	"github.com/lotas/docker-alerts/internal/handlers"
	"github.com/lotas/docker-alerts/internal/models"
	"github.com/lotas/docker-alerts/internal/notifications"
	"github.com/lotas/docker-alerts/internal/service"
)

func startApp(cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dockerClient, err := docker.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer dockerClient.Close()

	eventService := service.NewEventService(dockerClient)

	notifier := createNotifier(cfg)

	eventHandler := handlers.NewEventHandler(eventService, notifier)

	eventStream, err := eventService.StreamEvents(ctx)
	if err != nil {
		return fmt.Errorf("failed to start event stream: %w", err)
	}

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Event loop
	for {
		select {
		case event := <-eventStream.Events:
			if err := eventHandler.HandleContainerEvent(ctx, models.NewEventFromDocker(event)); err != nil {
				fmt.Printf("Error handling event: %v\n", err)
			}
		case err := <-eventStream.Errors:
			fmt.Printf("Error receiving event: %v\n", err)
		case <-sigChan:
			fmt.Println("Shutting down...")
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// createNotifier creates a notifier based on the configuration
func createNotifier(cfg *config.Config) notifications.Notifier {
	var notifiers []notifications.Notifier

	// Always add console notifier
	consoleNotifier := notifications.NewConsoleNotifier("DOCKER-ALERT",
		notifications.WithColor(),
		notifications.WithVerbose(),
	)
	notifiers = append(notifiers, consoleNotifier)

	// Add Slack notifier if configured
	if cfg.Notifications.Slack.WebhookURL != "" {
		slackNotifier := notifications.NewSlackNotifier(
			cfg.Notifications.Slack.WebhookURL,
			cfg.Notifications.Slack.Channel,
		)
		notifiers = append(notifiers, slackNotifier)
	}

	// Add Telegram notifier if configured
	if cfg.Notifications.Telegram.Token != "" && cfg.Notifications.Telegram.ChatID != "" {
		telegramNotifier := notifications.NewTelegramNotifier(
			cfg.Notifications.Telegram.Token,
			cfg.Notifications.Telegram.ChatID,
		)
		notifiers = append(notifiers, telegramNotifier)
	}

	// Add Email notifier if configured
	if cfg.Notifications.Email.SMTPHost != "" {
		emailNotifier := notifications.NewEmailNotifier(
			cfg.Notifications.Email.SMTPHost,
			cfg.Notifications.Email.SMTPPort,
			cfg.Notifications.Email.FromAddress,
			cfg.Notifications.Email.ToAddresses,
		)

		if cfg.Notifications.Email.SMTPUsername != "" && cfg.Notifications.Email.SMTPPassword != "" {
			emailNotifier.SetAuth(
				cfg.Notifications.Email.SMTPUsername,
				cfg.Notifications.Email.SMTPPassword,
			)
		}

		notifiers = append(notifiers, emailNotifier)
	}

	// Return multi-notifier if we have multiple notifiers
	if len(notifiers) > 1 {
		return notifications.NewMultiNotifier(notifiers...)
	}

	// Return console notifier as fallback
	return consoleNotifier
}
