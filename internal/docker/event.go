package docker

import (
	"fmt"

	"github.com/docker/docker/api/types/events"
	"github.com/lotas/docker-alerts/internal/notifications"
)

type EventActionMap map[string]map[string]bool
type ExitCodeMap map[string]string

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

var exitCodeMapping = ExitCodeMap{
	"0": "Success",
	// https://tldp.org/LDP/abs/html/exitcodes.html
	"1": "Application error",
	"2": "Misuse of builtin",
	// https://docs.docker.com/engine/containers/run/#exit-status
	"125": "Container failed to run",
	"126": "Container command cannot be invoked",
	"127": "Container command cannot be found",
	"128": "Invalid argument used on exit",

	// 128 + n Fatal
	"134": "Abnormal termination SIGABRT",
	"137": "Immediate termination SIGKILL",
	"139": "Segmentation Fault SIGSEGV",
	"143": "Graceful termination SIGTERM",

	"255": "Exit status out of range",
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
		if debug {
			fmt.Printf("Skipping unsupported event: %s:%s\n", e.Type, e.Action)
		}
		return false
	}

	return true
}

func getExitCodeDetails(exitCode string) string {
	if msg, ok := exitCodeMapping[exitCode]; ok {
		return msg
	}

	return ""
}

func (e Event) ToNotification() notifications.Notification {
	name, _ := e.Labels[containerNameLabel]
	project, _ := e.Labels[dockerComposeProjectLabel]
	service, _ := e.Labels[dockerComposeServiceLabel]
	exitCode, _ := e.Labels[exitCodeLabel]

	return notifications.Notification{
		Type:            e.Type,
		Action:          e.Action,
		Container:       e.Container,
		Image:           e.Image,
		Name:            name,
		Project:         project,
		Service:         service,
		ExitCode:        exitCode,
		ExitCodeDetails: getExitCodeDetails(exitCode),
	}
}
