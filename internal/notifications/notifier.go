package notifications

import (
	"context"
)

type Notification struct {
	Title   string
	Message string
	Level   string
}

type Notifier interface {
	Notify(ctx context.Context, notification Notification) error
}
