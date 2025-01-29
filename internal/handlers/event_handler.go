package handlers

import (
	"context"
	"github.com/lotas/docker-alerts/internal/models"
	"github.com/lotas/docker-alerts/internal/notifications"
	"github.com/lotas/docker-alerts/internal/service"
)

type EventHandler struct {
	eventService *service.EventService
	notifier     notifications.Notifier
}

func NewEventHandler(eventService *service.EventService, notifier notifications.Notifier) *EventHandler {
	return &EventHandler{
		eventService: eventService,
		notifier:     notifier,
	}
}

func (h *EventHandler) HandleContainerEvent(ctx context.Context, event models.Event) error {
	// Process container event
	if event.ShouldNotify() {
		notification := event.ToNotification()
		return h.notifier.Notify(ctx, notifications.Notification{
			Title:   notification.Title,
			Message: notification.Message,
			Level:   notification.Level,
		})
	}
	return nil
}
