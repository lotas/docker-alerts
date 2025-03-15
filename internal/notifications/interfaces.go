package notifications

import (
	"context"
)

type Notifier interface {
	Notify(ctx context.Context, event Event, debug bool) error
	NotifyMultiple(ctx context.Context, events []Event, debug bool) error
}
