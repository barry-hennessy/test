package redis_testcontainers

import (
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// NewFactory generates a sweet compatible DepFactory that  spins up Redis test
// containers of the given image.
func NewFactory(ctx context.Context, image string) func(t *testing.T) testcontainers.Container {
	return func(t *testing.T) testcontainers.Container {
		return NewContainer(t, ctx, image)
	}
}

// NewContainer sets up and runs a docker redis container for the given image.
//
// The container is cleaned up when the test ends.
func NewContainer(t *testing.T, ctx context.Context, image string) testcontainers.Container {
	listenWithinContainerOnPort := 6379
	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", listenWithinContainerOnPort)},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Error(err)
		return nil
	}

	t.Cleanup(func() {
		if err := redisC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	return redisC
}
