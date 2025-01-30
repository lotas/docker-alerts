package notifications

import (
	"context"
)

type MultiNotifier struct {
	notifiers []Notifier
}

func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{
		notifiers: notifiers,
	}
}

func (m *MultiNotifier) Notify(ctx context.Context, notification Notification) error {
	for _, notifier := range m.notifiers {
		if err := notifier.Notify(ctx, notification); err != nil {
			return err
		}
	}
	return nil
}
