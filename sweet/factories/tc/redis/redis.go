package redis

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisContainer struct {
	Image        string
	InternalPort string
}

func NewRedisContainer(image string) *RedisContainer {
	return &RedisContainer{
		Image:        image,
		InternalPort: "6379/tcp",
	}
}

func (rc *RedisContainer) Request() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        rc.Image,
		ExposedPorts: []string{rc.InternalPort},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}

}

func (rc *RedisContainer) Close(ctx context.Context, c testcontainers.Container) error {
	return c.Terminate(ctx)
}
