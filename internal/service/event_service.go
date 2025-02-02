package service

import (
	"context"
	"github.com/lotas/docker-alerts/internal/docker"
	"github.com/lotas/docker-alerts/internal/models"
	"github.com/lotas/docker-alerts/internal/notifications"
)

type EventService struct {
	dockerClient *docker.Client
	notifier notifications.Notifier
}

func NewEventService(dockerClient *docker.Client, notifier notifications.Notifier) *EventService {
	return &EventService{
		dockerClient: dockerClient,
		notifier: notifier,
	}
}

func (s *EventService) StreamEvents(ctx context.Context) (*docker.EventStream, error) {
	return s.dockerClient.StreamEvents(ctx)
}

func (s *EventService) HandleContainerEvent(ctx context.Context, event models.Event) error {
	if event.ShouldNotify() {
		notification := event.ToNotification()
		return s.notifier.Notify(ctx, notifications.Notification{
			Title:   notification.Title,
			Message: notification.Message,
			Level:   notification.Level,
		})
	}
	return nil
}
