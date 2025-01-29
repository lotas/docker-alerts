package service

import (
	"context"
	"github.com/lotas/docker-alerts/internal/docker"
)

type EventService struct {
	dockerClient *docker.Client
}

func NewEventService(dockerClient *docker.Client) *EventService {
	return &EventService{
		dockerClient: dockerClient,
	}
}

func (s *EventService) StreamEvents(ctx context.Context) (*docker.EventStream, error) {
	return s.dockerClient.StreamEvents(ctx)
}
