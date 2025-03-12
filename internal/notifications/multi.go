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

func (m *MultiNotifier) Notify(ctx context.Context, notification Notification, debug bool) error {
	for _, notifier := range m.notifiers {
		if err := notifier.Notify(ctx, notification, debug); err != nil {
			return err
		}
	}
	return nil
}

func (m *MultiNotifier) NotifyMultiple(ctx context.Context, notifications []Notification, debug bool) error {
	for _, notifier := range m.notifiers {
		if err := notifier.NotifyMultiple(ctx, notifications, debug); err != nil {
			return err
		}
	}
	return nil
}
