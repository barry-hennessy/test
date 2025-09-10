package elasticsearch

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ElasticsearchContainer struct {
	Image        string
	InternalPort string
}

func NewElasticsearchContainer(image string) *ElasticsearchContainer {
	return &ElasticsearchContainer{
		Image:        image,
		InternalPort: "9200/tcp",
	}
}

func (es *ElasticsearchContainer) Request() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image: es.Image,
		Env: map[string]string{
			"discovery.type":         "single-node",
			"xpack.security.enabled": "false",
		},
		ExposedPorts: []string{es.InternalPort},
		WaitingFor: wait.ForAll(
			wait.ForLog("started"),
			wait.ForHTTP("/").WithPort("9200/tcp"),
		),
	}
}

func (es *ElasticsearchContainer) Close(ctx context.Context, c testcontainers.Container) error {
	return c.Terminate(ctx)
}
