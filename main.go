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
	"github.com/lotas/docker-alerts/internal/models"
	"github.com/lotas/docker-alerts/internal/notifications"
	"github.com/lotas/docker-alerts/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
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

	info, err := dockerClient.Info(ctx)
	if err != nil {
	  return fmt.Errorf("failed to get Docker info: %w", err)
	}

	if cfg.Debug {
	 serverInfo := fmt.Sprintf("Docker version: %v\n", info.ServerVersion)
	 serverInfo += fmt.Sprintf("Docker host: %v\n", info.Name)
	 serverInfo += fmt.Sprintf("Type: %v\n", info.OSType)
	 serverInfo += fmt.Sprintf("Architecture: %v\n", info.Architecture)
	 serverInfo += fmt.Sprintf("CPUs: %v\n", info.NCPU)
	 serverInfo += fmt.Sprintf("Memory: %v MB\n", info.MemTotal / 1024 / 1024)
	 fmt.Println(serverInfo)
	}

	notifier := notifications.CreateNotifier(cfg)
	eventService := service.NewEventService(dockerClient, notifier)

	eventStream, err := eventService.StreamEvents(ctx)
	if err != nil {
		return fmt.Errorf("failed to start event stream: %w", err)
	}

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case event := <-eventStream.Events:
			if err := eventService.HandleContainerEvent(ctx, models.NewEventFromDocker(event)); err != nil {
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
