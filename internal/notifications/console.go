package notifications

import (
	"context"
	"fmt"
	"log"
)

type ConsoleNotifier struct {
	prefix  string
	colored bool
}

type ConsoleOption func(*ConsoleNotifier)

func WithColor() ConsoleOption {
	return func(n *ConsoleNotifier) {
		n.colored = true
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

func (c *ConsoleNotifier) Notify(ctx context.Context, event Event, debug bool) error {
	var message string
	if c.colored {
		message = fmt.Sprintf(Yellow+"[%s]"+Reset+" "+Blue+"%s"+Reset+"%s",
			c.prefix,
			event.Type,
			event.ANSI(),
		)
	} else {
		message = fmt.Sprintf("[%s] %s%s",
			c.prefix,
			event.Type,
			event.Text(),
		)
	}

	log.Print(message)
	return nil
}

func (c *ConsoleNotifier) NotifyMultiple(ctx context.Context, events []Event, debug bool) error {
	for _, n := range events {
		c.Notify(ctx, n, debug)
	}
	return nil
}
