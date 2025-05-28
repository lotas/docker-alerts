package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/system"
)

type DockerAPIClient interface {
	Info(ctx context.Context) (system.Info, error)
	Events(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error)
	Close() error
}
