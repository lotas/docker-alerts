package docker

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/system"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockDockerClient is a custom struct implementing necessary methods for testing
type mockDockerClient struct {
	infoFunc   func(ctx context.Context) (system.Info, error)
	eventsFunc func(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error)
}

func (m *mockDockerClient) Info(ctx context.Context) (system.Info, error) {
	return m.infoFunc(ctx)
}

func (m *mockDockerClient) Events(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error) {
	return m.eventsFunc(ctx, options)
}

func (m *mockDockerClient) Close() error {
	return nil
}

func TestInfo(t *testing.T) {
	t.Run("successful info retrieval", func(t *testing.T) {
		// Create mock info data
		mockInfo := system.Info{
			ServerVersion: "20.10.12",
			Name:          "docker-test-host",
			OSType:        "linux",
			Architecture:  "x86_64",
			NCPU:          4,
			MemTotal:      8589934592, // 8GB in bytes
		}

		// Create a mock docker client
		mockClient := &mockDockerClient{
			infoFunc: func(ctx context.Context) (system.Info, error) {
				return mockInfo, nil
			},
		}

		// Create our Docker client with the mocked client
		c := &Client{
			cli: mockClient,
		}

		// Call the Info method
		ctx := context.Background()
		info, serverInfo, err := c.Info(ctx)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, mockInfo, info)

		// Check serverInfo string contains expected information
		expectedParts := []string{
			"Docker version: 20.10.12",
			"Docker host: docker-test-host",
			"Type: linux",
			"Architecture: x86_64",
			"CPUs: 4",
			"Memory: 8192 MB", // 8GB in MB
		}

		for _, part := range expectedParts {
			assert.Contains(t, serverInfo, part)
		}
	})

	t.Run("error retrieving info", func(t *testing.T) {
		// Create a mock docker client that returns an error
		mockClient := &mockDockerClient{
			infoFunc: func(ctx context.Context) (system.Info, error) {
				return system.Info{}, fmt.Errorf("connection error")
			},
		}

		// Create our Docker client with the mocked client
		c := &Client{
			cli: mockClient,
		}

		// Call the Info method
		ctx := context.Background()
		_, serverInfo, err := c.Info(ctx)

		// Assertions
		require.Error(t, err)
		assert.Equal(t, "", serverInfo)
		assert.Contains(t, err.Error(), "connection error")
	})
}
