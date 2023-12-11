package nats

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type NatsContainer struct {
	Image          string
	ClientPort     string
	ClusterPort    string
	MonitoringPort string
}

func NewNatsContainer(image string) *NatsContainer {
	return &NatsContainer{
		Image:          image,
		ClientPort:     "4222/tcp",
		ClusterPort:    "6222/tcp",
		MonitoringPort: "8222/tcp",
	}
}

func (nc *NatsContainer) Request() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        nc.Image,
		Cmd:          []string{"--js"},
		ExposedPorts: []string{nc.ClientPort, nc.ClusterPort, nc.MonitoringPort},
		WaitingFor:   wait.ForLog(".*Server is ready").AsRegexp(),
	}
}

func (nc *NatsContainer) Close(ctx context.Context, c testcontainers.Container) error {
	return c.Terminate(ctx)
}
