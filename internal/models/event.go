package models

import (
	"fmt"

	"github.com/docker/docker/api/types/events"
	"github.com/lotas/docker-alerts/internal/notifications"
)

type Event struct {
	Type      string
	Action    string
	Container string
	Image     string
	Time      int64
	Status    string
	Labels    map[string]string
}

func NewEventFromDocker(msg events.Message) Event {
	return Event{
		Type:      string(msg.Type),
		Action:    string(msg.Action),
		Container: msg.Actor.ID,
		Image:     msg.Actor.Attributes["image"],
		Time:      msg.Time,
		Status:    msg.Status,
		Labels:    msg.Actor.Attributes,
	}
}

func (e Event) ShouldNotify(debug bool) bool {
	if debug {
		fmt.Printf("ShouldNotify type: %+v action: %+v", e.Type, e.Action)
	}
	return true
}

func (e Event) ToNotification() notifications.Notification {
	return notifications.Notification{
		Title:   fmt.Sprintf("%s event: %s", e.Type, e.Action),
		Message: fmt.Sprintf("Container: %s\nImage: %s", e.Container, e.Image),
	}
}
