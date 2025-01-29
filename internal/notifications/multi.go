package notifications

import (
	"context"
)

// MultiNotifier allows sending notifications to multiple notifiers
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier creates a new MultiNotifier
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{
		notifiers: notifiers,
	}
}

// Notify sends the notification to all configured notifiers
func (m *MultiNotifier) Notify(ctx context.Context, notification Notification) error {
	for _, notifier := range m.notifiers {
		if err := notifier.Notify(ctx, notification); err != nil {
			return err
		}
	}
	return nil
}
