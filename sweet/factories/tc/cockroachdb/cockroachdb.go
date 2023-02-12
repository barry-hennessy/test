package cockroachdb

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type CockroachDBContainer struct {
	InternalPorts []string
	Image         string
}

func NewCockroachDBContainer(image string) *CockroachDBContainer {
	return &CockroachDBContainer{
		Image:         image,
		InternalPorts: []string{"26257/tcp", "8080/tcp"},
	}
}

func (crdb *CockroachDBContainer) Request() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        crdb.Image,
		ExposedPorts: crdb.InternalPorts,
		WaitingFor:   wait.ForHTTP("/health").WithPort("8080"),
		Cmd:          []string{"start-single-node", "--insecure"},
	}
}

func (crdb *CockroachDBContainer) Close(ctx context.Context, c testcontainers.Container) error {
	return c.Terminate(ctx)
}
