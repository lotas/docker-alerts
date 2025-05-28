package notifications

import (
	"bytes"
	"context"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureOutput captures log output for testing.
// It returns a function that should be called (e.g., with defer) to restore the original log output
// and return the captured string.
func captureOutput(t *testing.T) (*bytes.Buffer, func() string) {
	t.Helper()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0) // Remove timestamp and other prefixes for predictable output

	// It's important to guard access to the buffer if tests run in parallel,
	// though for these specific tests, they modify a global (log.SetOutput)
	// so parallel execution of tests that use this function is problematic.
	// For simplicity, assuming tests using this are run serially or one at a time.
	// If parallel tests are needed, each test would need its own logger instance.

	return &buf, func() string {
		log.SetOutput(os.Stderr) // Restore default logger
		return buf.String()
	}
}

func TestNewConsoleNotifier(t *testing.T) {
	t.Run("default prefix", func(t *testing.T) {
		notifier := NewConsoleNotifier("")
		assert.Equal(t, "DOCKER-EVENT", notifier.prefix, "Expected default prefix")
		assert.False(t, notifier.colored, "Expected colored to be false by default")
	})

	t.Run("custom prefix", func(t *testing.T) {
		notifier := NewConsoleNotifier("CUSTOM_PREFIX")
		assert.Equal(t, "CUSTOM_PREFIX", notifier.prefix, "Expected custom prefix")
		assert.False(t, notifier.colored, "Expected colored to be false")
	})

	t.Run("with color option", func(t *testing.T) {
		notifier := NewConsoleNotifier("", WithColor())
		assert.True(t, notifier.colored, "Expected colored to be true with WithColor option")
	})

	t.Run("custom prefix and with color option", func(t *testing.T) {
		notifier := NewConsoleNotifier("MY_EVENTS", WithColor())
		assert.Equal(t, "MY_EVENTS", notifier.prefix, "Expected custom prefix")
		assert.True(t, notifier.colored, "Expected colored to be true")
	})
}

func TestConsoleNotifier_Notify(t *testing.T) {
	sampleEvent := Event{
		Type:   "container",
		Action: "start",
		Name:   "test-container",
		Image:  "test-image:latest",
	}

	t.Run("notify without color", func(t *testing.T) {
		_, cleanup := captureOutput(t)
		defer cleanup()

		notifier := NewConsoleNotifier("TEST")
		err := notifier.Notify(context.Background(), sampleEvent, false)
		require.NoError(t, err)

		output := cleanup()
		expectedLog := "[TEST] container start test-container (test-image:latest)"
		// Trim newline which log.Print adds
		assert.Equal(t, expectedLog, strings.TrimSpace(output))
	})

	t.Run("notify with color", func(t *testing.T) {
		_, cleanup := captureOutput(t)
		defer cleanup()

		notifier := NewConsoleNotifier("COLORTEST", WithColor())
		err := notifier.Notify(context.Background(), sampleEvent, false)
		require.NoError(t, err)

		output := cleanup()
		// Check for prefix and parts of the ANSI output. Exact ANSI codes can be brittle.
		assert.Contains(t, output, "[COLORTEST]")
		assert.Contains(t, output, sampleEvent.Name) // Name should be in the output
		assert.Contains(t, output, Yellow)           // Action should be yellow
		assert.Contains(t, output, Cyan)             // Name should be cyan
		assert.Contains(t, output, Green)            // Image should be green
		assert.Contains(t, output, Reset)            // Reset code should be present
	})

	t.Run("notify with custom message", func(t *testing.T) {
		_, cleanup := captureOutput(t)
		defer cleanup()

		customMessageEvent := Event{
			Message: "This is a custom message",
		}
		notifier := NewConsoleNotifier("CUSTOM")
		err := notifier.Notify(context.Background(), customMessageEvent, false)
		require.NoError(t, err)
		output := cleanup()
		expectedLog := "[CUSTOM] This is a custom message"
		assert.Equal(t, expectedLog, strings.TrimSpace(output))
	})
}

// TestConsoleNotifier_NotifyMultiple tests the NotifyMultiple method
func TestConsoleNotifier_NotifyMultiple(t *testing.T) {
	event1 := Event{Type: "container", Action: "start", Name: "c1", Image: "img1"}
	event2 := Event{Type: "container", Action: "stop", Name: "c2", Image: "img2"}

	t.Run("notify multiple without color", func(t *testing.T) {
		_, cleanup := captureOutput(t)
		defer cleanup()

		notifier := NewConsoleNotifier("MULTI")
		err := notifier.NotifyMultiple(context.Background(), []Event{event1, event2}, false)
		require.NoError(t, err)

		output := cleanup()
		// log.Print adds a newline for each call
		expectedLogs := "[MULTI] container start c1 (img1)\n[MULTI] container stop c2 (img2)"
		assert.Equal(t, expectedLogs, strings.TrimSpace(output))
	})

	t.Run("notify multiple with color", func(t *testing.T) {
		_, cleanup := captureOutput(t)
		defer cleanup()

		notifier := NewConsoleNotifier("COLORMULTI", WithColor())
		err := notifier.NotifyMultiple(context.Background(), []Event{event1, event2}, false)
		require.NoError(t, err)

		output := cleanup()
		assert.Contains(t, output, "[COLORMULTI]")
		assert.Contains(t, output, "c1")
		assert.Contains(t, output, "c2")
		assert.Contains(t, output, Yellow) // Check for color codes
		assert.Contains(t, output, Reset)
	})

	t.Run("notify multiple with empty slice", func(t *testing.T) {
		_, cleanup := captureOutput(t)
		defer cleanup()

		notifier := NewConsoleNotifier("EMPTY")
		err := notifier.NotifyMultiple(context.Background(), []Event{}, false)
		require.NoError(t, err)

		output := cleanup()
		assert.Equal(t, "", strings.TrimSpace(output), "Expected no log output for empty event slice")
	})
}

// Note: The init() function in notification.go which parses templates
// might panic if templates are invalid. If these tests fail,
// ensure that the templates in notification.go are correct.
// However, console_notifier specifically uses event.Text() and event.ANSI(),
// so direct template issues from notification.go's init() might not surface here
// unless those methods fail.
