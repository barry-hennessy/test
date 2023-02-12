package mongodb

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type MongoDBContainer struct {
	Image        string
	InternalPort string
}

func NewMongoDBContainer(image string) *MongoDBContainer {
	return &MongoDBContainer{
		Image:        image,
		InternalPort: "27017/tcp",
	}
}

func (rc *MongoDBContainer) Request() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        rc.Image,
		ExposedPorts: []string{rc.InternalPort},
		WaitingFor: wait.ForAll(
			wait.ForLog("Waiting for connections"),
			wait.ForListeningPort("27017/tcp"),
		),
	}

}

func (rc *MongoDBContainer) Close(ctx context.Context, c testcontainers.Container) error {
	return c.Terminate(ctx)
}
