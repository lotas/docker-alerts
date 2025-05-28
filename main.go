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
	"github.com/lotas/docker-alerts/internal/notifications"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Debug {
		cfg.PrintValues()
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

	_, infoStr, err := dockerClient.Info(ctx)
	if err != nil {
		return fmt.Errorf("failed to get Docker info: %w", err)
	}

	notifier := notifications.CreateNotifier(cfg)

	err0 := notifier.Notify(ctx, notifications.Event{
		Type:    "Server info",
		Message: infoStr,
	}, cfg.Debug)
	if err0 != nil {
		return fmt.Errorf("Failed to notify: %w", err0)
	}

	eventStream, err := dockerClient.StreamEvents(ctx)
	if err != nil {
		return fmt.Errorf("failed to start event stream: %w", err)
	}

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case event := <-eventStream.Events:
			evt := notifications.NewEventFromDocker(event)
			if evt.ShouldNotify(cfg.Debug) {
				err := notifier.Notify(ctx, evt, cfg.Debug)
				if err != nil {
					fmt.Printf("Error sending event %+v", err)
				}
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
