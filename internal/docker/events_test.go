package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockDockerEventsClient is a helper struct to mock the DockerAPIClient for events_test.go
type mockDockerEventsClient struct {
	mockDockerClient // Embed the mock from client_test.go if common methods are needed
	eventsFunc       func(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error)
}

// Overriding the Events method from the embedded mockDockerClient if necessary,
// or just implementing it if mockDockerClient doesn't have it.
func (m *mockDockerEventsClient) Events(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error) {
	if m.eventsFunc != nil {
		return m.eventsFunc(ctx, options)
	}
	// Default behavior if not set by a specific test
	msgChan := make(chan events.Message)
	errChan := make(chan error)
	// Close channels immediately for tests not expecting events.
	close(msgChan)
	close(errChan)
	return msgChan, errChan
}

// Ensure mockDockerEventsClient implements DockerAPIClient.
// We might need to add other methods from DockerAPIClient if StreamEvents or its setup calls them.
// For now, assuming only Events() and Close() (from embedded mockDockerClient) are relevant.
var _ DockerAPIClient = (*mockDockerEventsClient)(nil)

func TestStreamEvents(t *testing.T) {
	t.Run("successful event stream without filters", func(t *testing.T) {
		actualMockedEventsChan := make(chan events.Message)
		actualMockedErrorsChan := make(chan error)
		defer close(actualMockedEventsChan)
		defer close(actualMockedErrorsChan)

		var expectedEventsChan <-chan events.Message = actualMockedEventsChan
		var expectedErrorsChan <-chan error = actualMockedErrorsChan

		mockAPIClient := &mockDockerEventsClient{
			eventsFunc: func(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error) {
				// Check that filters are empty
				assert.Equal(t, 0, options.Filters.Len(), "Expected no filters to be set")
				return actualMockedEventsChan, actualMockedErrorsChan
			},
		}

		// Create our Docker client with the mocked API client
		c := &Client{cli: mockAPIClient}

		eventStream, err := c.StreamEvents(context.Background())
		require.NoError(t, err)
		require.NotNil(t, eventStream)

		assert.Equal(t, expectedEventsChan, eventStream.Events, "Events channel does not match")
		assert.Equal(t, expectedErrorsChan, eventStream.Errors, "Errors channel does not match")
	})

	t.Run("successful event stream with filters", func(t *testing.T) {
		actualMockedEventsChan := make(chan events.Message)
		actualMockedErrorsChan := make(chan error)
		defer close(actualMockedEventsChan)
		defer close(actualMockedErrorsChan)

		var expectedEventsChan <-chan events.Message = actualMockedEventsChan
		var expectedErrorsChan <-chan error = actualMockedErrorsChan

		expectedFilters := filters.NewArgs()
		expectedFilters.Add("type", "container")

		mockAPIClient := &mockDockerEventsClient{
			eventsFunc: func(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error) {
				// Check that filters are set as expected
				assert.True(t, options.Filters.ExactMatch("type", "container"), "Expected 'type' filter to be 'container'")
				return actualMockedEventsChan, actualMockedErrorsChan
			},
		}

		c := &Client{cli: mockAPIClient}

		eventStream, err := c.StreamEvents(context.Background(), expectedFilters)
		require.NoError(t, err)
		require.NotNil(t, eventStream)

		assert.Equal(t, expectedEventsChan, eventStream.Events, "Events channel does not match")
		assert.Equal(t, expectedErrorsChan, eventStream.Errors, "Errors channel does not match")
	})

	t.Run("event stream with multiple filter arguments", func(t *testing.T) {
		// This test verifies that only the first filterArgs is used, as per implementation.
		actualMockedEventsChan := make(chan events.Message)
		actualMockedErrorsChan := make(chan error)
		defer close(actualMockedEventsChan)
		defer close(actualMockedErrorsChan)
		// Not asserting channel equality here as the main point is to check filter passing.
		// var expectedEventsChan <-chan events.Message = actualMockedEventsChan
		// var expectedErrorsChan <-chan error = actualMockedErrorsChan

		filters1 := filters.NewArgs()
		filters1.Add("event", "start")

		filters2 := filters.NewArgs() // This should be ignored by StreamEvents
		filters2.Add("event", "stop")

		mockAPIClient := &mockDockerEventsClient{
			eventsFunc: func(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error) {
				assert.True(t, options.Filters.ExactMatch("event", "start"), "Expected 'event' filter to be 'start'")
				assert.False(t, options.Filters.ExactMatch("event", "stop"), "Filter 'stop' should not be present")
				return actualMockedEventsChan, actualMockedErrorsChan
			},
		}
		c := &Client{cli: mockAPIClient}
		_, err := c.StreamEvents(context.Background(), filters1, filters2)
		require.NoError(t, err)

	})
}

// It might be necessary to define a Close() method on mockDockerEventsClient
// if the Client's Close() method is ever called, or if other parts of the system expect it.
// If mockDockerClient (embedded) already provides a no-op Close(), this is fine.
// func (m *mockDockerEventsClient) Close() error { return nil }
