package notifications

import (
	"context"
	"fmt"
	"log"
	"time"
)

type ConsoleNotifier struct {
	prefix  string
	colored bool
	verbose bool
}

type ConsoleOption func(*ConsoleNotifier)

func WithColor() ConsoleOption {
	return func(n *ConsoleNotifier) {
		n.colored = true
	}
}

func WithVerbose() ConsoleOption {
	return func(n *ConsoleNotifier) {
		n.verbose = true
	}
}

func NewConsoleNotifier(prefix string, opts ...ConsoleOption) *ConsoleNotifier {
	if prefix == "" {
		prefix = "DOCKER-EVENT"
	}

	n := &ConsoleNotifier{
		prefix: prefix,
	}

	for _, opt := range opts {
		opt(n)
	}

	return n
}

func (c *ConsoleNotifier) Notify(ctx context.Context, notification Notification) error {
	timestamp := time.Now().Format(time.RFC3339)

	var message string
	if c.colored {
		// Add ANSI color codes
		message = fmt.Sprintf("\033[1;34m[%s]\033[0m \033[1;32m[%s]\033[0m \033[1;33m%s\033[0m\n%s\n",
			c.prefix,
			timestamp,
			notification.Title,
			notification.Message,
		)
	} else {
		message = fmt.Sprintf("[%s] [%s] %s\n%s\n",
			c.prefix,
			timestamp,
			notification.Title,
			notification.Message,
		)
	}

	if c.verbose && notification.Level != "" {
		message += fmt.Sprintf("Level: %s\n", notification.Level)
	}

	log.Print(message)
	return nil
}
