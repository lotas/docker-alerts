package notifications

import (
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/events"
)

func TestNewEventFromDocker(t *testing.T) {
	// Create a mock Docker event
	msg := events.Message{
		Type:   "container",
		Action: "start",
		Actor: events.Actor{
			ID: "abc123",
			Attributes: map[string]string{
				"image":                      "nginx:latest",
				"name":                       "web-server",
				"com.docker.compose.project": "myproject",
				"com.docker.compose.service": "webapp",
				"exitCode":                   "0",
				"execDuration":               "10",
			},
		},
		Time:   time.Now().Unix(),
		Status: "running",
	}

	event := NewEventFromDocker(msg)

	// Verify event fields are correctly populated
	if event.Type != "container" {
		t.Errorf("Expected type 'container', got %s", event.Type)
	}
	if event.Action != "start" {
		t.Errorf("Expected action 'start', got %s", event.Action)
	}
	if event.Container != "abc123" {
		t.Errorf("Expected container ID 'abc123', got %s", event.Container)
	}
	if event.Image != "nginx:latest" {
		t.Errorf("Expected image 'nginx:latest', got %s", event.Image)
	}
	if event.Name != "web-server" {
		t.Errorf("Expected name 'web-server', got %s", event.Name)
	}
	if event.Project != "myproject" {
		t.Errorf("Expected project 'myproject', got %s", event.Project)
	}
	if event.Service != "webapp" {
		t.Errorf("Expected service 'webapp', got %s", event.Service)
	}
	if event.ExitCode != "0" {
		t.Errorf("Expected exit code '0', got %s", event.ExitCode)
	}
	if event.ExitCodeDetails != "Success" {
		t.Errorf("Expected exit code details 'Success', got %s", event.ExitCodeDetails)
	}
	if event.ExecDuration != "10" {
		t.Errorf("Expected execution duration '10', got %s", event.ExecDuration)
	}
}

func TestShouldNotify(t *testing.T) {
	tests := []struct {
		name     string
		event    Event
		expected bool
	}{
		{
			name: "Container start event",
			event: Event{
				Type:   "container",
				Action: "start",
			},
			expected: true,
		},
		{
			name: "Container die event",
			event: Event{
				Type:   "container",
				Action: "die",
			},
			expected: true,
		},
		{
			name: "Container healthy event",
			event: Event{
				Type:   "container",
				Action: "health_status: healthy",
			},
			expected: true,
		},
		{
			name: "Container unhealthy event",
			event: Event{
				Type:   "container",
				Action: "health_status: unhealthy",
			},
			expected: true,
		},
		{
			name: "Connection message event",
			event: Event{
				Type:   "connection",
				Action: "message",
			},
			expected: true,
		},
		{
			name: "Unsupported event type",
			event: Event{
				Type:   "network",
				Action: "connect",
			},
			expected: false,
		},
		{
			name: "Unsupported action for container",
			event: Event{
				Type:   "container",
				Action: "create",
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.event.ShouldNotify(false)
			if result != tc.expected {
				t.Errorf("Expected ShouldNotify to return %v but got %v", tc.expected, result)
			}
		})
	}
}

func TestEventText(t *testing.T) {
	event := Event{
		Type:         "container",
		Action:       "start",
		Name:         "web-server",
		Image:        "nginx:latest",
		Project:      "myproject",
		Service:      "webapp",
		ExitCode:     "0",
		ExecDuration: "10",
	}

	result := event.Text()
	expectedParts := []string{"container", "start", "web-server", "nginx:latest"}
	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("Expected text to contain '%s', but got: %s", part, result)
		}
	}

	// Test with custom message
	event.Message = "Custom notification message"
	result = event.Text()
	if result != "Custom notification message" {
		t.Errorf("Expected custom message, got: %s", result)
	}
}

func TestEventMarkdown(t *testing.T) {
	event := Event{
		Type:     "container",
		Action:   "die",
		Name:     "web-server",
		Image:    "nginx:latest",
		Project:  "myproject",
		Service:  "webapp",
		ExitCode: "1",
	}

	result := event.Markdown()
	expectedParts := []string{"container", "*stop*", "`web-server`", "`nginx:latest`"}
	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("Expected markdown to contain '%s', but got: %s", part, result)
		}
	}
}

func TestEventANSI(t *testing.T) {
	event := Event{
		Type:    "container",
		Action:  "health_status: healthy",
		Name:    "web-server",
		Image:   "nginx:latest",
		Project: "myproject",
		Service: "webapp",
	}

	result := event.ANSI()
	expectedParts := []string{"container", "healthy", "web-server", "nginx"}
	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("Expected ANSI text to contain '%s', but got: %s", part, result)
		}
	}
}

func TestExitCodeDetails(t *testing.T) {
	tests := []struct {
		exitCode string
		expected string
	}{
		{"0", "Success"},
		{"1", "Application error"},
		{"137", "Immediate termination SIGKILL"},
		{"143", "Graceful termination SIGTERM"},
		{"999", ""}, // Unknown exit code
	}

	for _, tc := range tests {
		t.Run("ExitCode_"+tc.exitCode, func(t *testing.T) {
			details := getExitCodeDetails(tc.exitCode)
			if details != tc.expected {
				t.Errorf("Expected exit code %s to map to '%s', got '%s'",
					tc.exitCode, tc.expected, details)
			}
		})
	}
}

func TestActionNameMapping(t *testing.T) {
	// We can test the ActionName template function by using the templates directly
	event := Event{
		Type:   "container",
		Action: "die",
		Name:   "test-container",
		Image:  "test-image",
	}

	result := event.Text()
	if !strings.Contains(result, "stop") {
		t.Errorf("Expected 'die' action to be mapped to 'stop', but got: %s", result)
	}

	event.Action = "health_status: healthy"
	result = event.Text()
	if !strings.Contains(result, "healthy") {
		t.Errorf("Expected health action to be mapped to 'healthy', but got: %s", result)
	}
}
