package redisfactory

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/testcontainers/testcontainers-go"

	redis_testcontainer_factory "github.com/barry-hennessy/test/sweet/factories/testcontainers/redis"
	"github.com/docker/go-connections/nat"
)

// NewContainerClientFactory creates a sweet compatible DepFactory that runs a
// redis docker container and returns a redis client connected to it.
//
// The container and client are cleaned up once the test ends.
func NewContainerClientFactory(ctx context.Context, image string, opts *redis.Options) func(t *testing.T) *redis.Client {
	return func(t *testing.T) *redis.Client {
		c := redis_testcontainer_factory.NewFactory(ctx, image)(t)
		return NewContainerClient(t, ctx, c, opts)
	}
}

// NewContainerClient creates a new redis client listening to a redis server
// running in the test container given.
//
// Redis clients connect using the options given; `Addr` will be overridden with
// the correct host and port for the created test container.
func NewContainerClient(t *testing.T, ctx context.Context, c testcontainers.Container, opts *redis.Options) *redis.Client {
	containerInternalRedisPort := "6379"
	mappedPort, err := c.MappedPort(ctx, nat.Port(containerInternalRedisPort))
	if err != nil {
		t.Error(err)
		return nil
	}

	hostIP, err := c.Host(ctx)
	if err != nil {
		t.Error(err)
		return nil
	}

	opts.Addr = fmt.Sprintf("%s:%s", hostIP, mappedPort.Port())

	return NewClient(t, ctx, opts)
}

// NewClientFactory creates a sweet compatible DepFactory that creates fresh
// redis clients with the given connection options.
func NewClientFactory(ctx context.Context, opts *redis.Options) func(t *testing.T) *redis.Client {
	return func(t *testing.T) *redis.Client {
		return NewClient(t, ctx, opts)
	}
}

// NewClient creates a new redis client with the given connection options.
func NewClient(t *testing.T, ctx context.Context, opts *redis.Options) *redis.Client {
	client := redis.NewClient(opts).WithContext(ctx)
	t.Cleanup(func() {
		client.Close()
	})
	return client
}
