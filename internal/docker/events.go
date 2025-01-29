package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
)

type EventStream struct {
	Events <-chan events.Message
	Errors <-chan error
}

func (c *Client) StreamEvents(ctx context.Context, filterArgs ...filters.Args) (*EventStream, error) {
	var opts types.EventsOptions
	if len(filterArgs) > 0 {
		opts.Filters = filterArgs[0]
	}

	eventsChan, errorsChan := c.cli.Events(ctx, opts)

	return &EventStream{
		Events: eventsChan,
		Errors: errorsChan,
	}, nil
}
