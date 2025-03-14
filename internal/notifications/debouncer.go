package notifications

import (
	"context"
	"sync"
	"time"
)

type DebouncerNotifier struct {
	notifier    Notifier
	events      []Event
	ctx         context.Context
	debug       bool
	lastSent    time.Time
	timer       *time.Timer
	mu          sync.Mutex
	minInterval time.Duration
	isScheduled bool
}

func NewDebouncerNotifier(notifier Notifier, minInterval time.Duration) *DebouncerNotifier {
	if minInterval == 0 {
		minInterval = 5 * time.Second
	}

	return &DebouncerNotifier{
		notifier:    notifier,
		lastSent:    time.Time{},
		minInterval: minInterval,
	}
}

func (d *DebouncerNotifier) NotifyMultiple(ctx context.Context, events []Event, debug bool) error {
	// shouldn't be really called but ok
	for _, n := range events {
		d.Notify(ctx, n, debug)
	}
	return nil
}

func (d *DebouncerNotifier) Notify(ctx context.Context, n Event, debug bool) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.events = append(d.events, n)
	d.ctx = ctx
	d.debug = debug

	timeElapsed := time.Since(d.lastSent)
	if timeElapsed >= d.minInterval {
		d.sendAllLocked()
		return nil
	}

	if !d.isScheduled {
		if d.timer != nil {
			d.timer.Stop()
		}

		delay := d.minInterval - timeElapsed

		d.timer = time.AfterFunc(delay, func() {
			d.mu.Lock()
			defer d.mu.Unlock()

			d.isScheduled = false
			d.sendAllLocked()
		})

		d.isScheduled = true
	}

	return nil
}

// must be called when lock is held
func (d *DebouncerNotifier) sendAllLocked() {
	if len(d.events) == 0 {
		return
	}

	d.notifier.NotifyMultiple(d.ctx, d.events, d.debug)

	d.lastSent = time.Now()
	d.events = nil
}

func (d *DebouncerNotifier) Close() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
}
