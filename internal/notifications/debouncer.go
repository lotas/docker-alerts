package notifications

import (
	"context"
)

type DebouncerNotifier struct {
	notifier     Notifier
	notification Notification
}

func NewDebouncerNotifier(notifier Notifier) *DebouncerNotifier {
	return &DebouncerNotifier{
		notifier: notifier,
	}
}

func (d *DebouncerNotifier) Notify(ctx context.Context, n Notification) error {
	if n.IsSame(d.notification) {
		d.notification.TimesSeen += 1
	} else {
		d.notification = n
	}
	// add more logic to avoid sending same messages too frequently

	return d.notifier.Notify(ctx, d.notification)
}
