package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/events"
	"github.com/lotas/docker-alerts/internal/notifications"
)

type EventActionMap map[string]map[string]bool

const containerNameLabel = "name"
const dockerComposeProjectLabel = "com.docker.compose.project"
const dockerComposeServiceLabel = "com.docker.compose.service"
const execDurationLabel = "execDuration"
const exitCodeLabel = "exitCode"

var SupportedEvents = EventActionMap{
	"container": {
		"start":                    true,
		"die":                      true,
		"health_status: healthy":   true,
		"health_status: unhealthy": true,
	},
	"connection": {
		"message": true,
	},
}

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
	supported := false
	if actionMap, ok := SupportedEvents[e.Type]; ok {
		supported, _ = actionMap[e.Action]
	}
	if !supported {
		return false
	}

	if debug {
		jsonStr, _ := json.MarshalIndent(e, "", "  ")
		fmt.Printf("Should notify:\n%v\n", string(jsonStr))
	}
	return true
}

func (e Event) ToNotification() notifications.Notification {
	info := []string{
		fmt.Sprintf("Container: `%s`", e.Container[0:20]),
		fmt.Sprintf("Image: `%s`", e.Image),
	}
	if name, ok := e.Labels[containerNameLabel]; ok {
		info = append(info, fmt.Sprintf("Name: `%s`", name))
	}
	if project, ok := e.Labels[dockerComposeProjectLabel]; ok {
		info = append(info, fmt.Sprintf("Project: `%s`", project))
	}
	if service, ok := e.Labels[dockerComposeServiceLabel]; ok {
		info = append(info, fmt.Sprintf("Service: `%s`", service))
	}
	if exitCode, ok := e.Labels[exitCodeLabel]; ok {
		info = append(info, fmt.Sprintf("Exit code: `%v`", exitCode))
	}

	return notifications.Notification{
		Title:   fmt.Sprintf("%s %s", e.Type, e.Action),
		Message: strings.Join(info, "\n"),
	}
}
