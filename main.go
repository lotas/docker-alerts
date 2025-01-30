package main

import (
	"fmt"
	"log"

	"context"
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

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Debug {
		fmt.Printf("Configuration:\n%+v\n", cfg)
	}

	if err := startApp(cfg); err != nil {
		log.Fatal(err)
	}
}

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

	// graceful shutdown
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

func createNotifier(cfg *config.Config) notifications.Notifier {
	var notifiers []notifications.Notifier

	consoleNotifier := notifications.NewConsoleNotifier("DOCKER-ALERT",
		notifications.WithColor(),
		notifications.WithVerbose(),
	)
	notifiers = append(notifiers, consoleNotifier)

	if cfg.SlackWebhookURL != "" {
		slackNotifier := notifications.NewSlackNotifier(
			cfg.SlackWebhookURL,
			cfg.SlackChannel,
		)
		notifiers = append(notifiers, slackNotifier)
	}

	if cfg.TelegramToken != "" && cfg.TelegramChatID != "" {
		telegramNotifier := notifications.NewTelegramNotifier(
			cfg.TelegramToken,
			cfg.TelegramChatID,
		)
		notifiers = append(notifiers, telegramNotifier)
	}

	if cfg.EmailSMTPHost != "" {
		emailNotifier := notifications.NewEmailNotifier(
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
		return notifications.NewMultiNotifier(notifiers...)
	}

	return consoleNotifier
}
