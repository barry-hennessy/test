package tc

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

type Container interface {
	Request() testcontainers.ContainerRequest
	Close(ctx context.Context, container testcontainers.Container) error
}

// NewContainer sets up and runs a docker container for the given image.
//
// The container is cleaned up when the test ends.
func NewContainer(t *testing.T, ctx context.Context, c Container) testcontainers.Container {
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: c.Request(),
		Started:          true,
		Logger:           testcontainers.TestLogger(t),
	})

	if err != nil {
		t.Errorf("could not start container: %s", err)
		return nil
	}

	t.Cleanup(func() {
		if err := c.Close(ctx, container); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})

	return container
}

// NewFactory generates a sweet compatible DepFactory that  spins up Redis test
// containers of the given image.
func NewFactory(ctx context.Context, c Container) func(t *testing.T) testcontainers.Container {
	return func(t *testing.T) testcontainers.Container {
		return NewContainer(t, ctx, c)
	}
}
