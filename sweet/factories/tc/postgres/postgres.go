package postgres

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	Image        string
	InternalPort string
}

func NewPostgresContainer(image string) *PostgresContainer {
	return &PostgresContainer{
		Image:        image,
		InternalPort: "5432/tcp",
	}
}

func (pg *PostgresContainer) Request() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image: pg.Image,
		Env: map[string]string{
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_DB":       "postgres",
		},
		ExposedPorts: []string{pg.InternalPort},
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
	}
}

func (pc *PostgresContainer) Close(ctx context.Context, c testcontainers.Container) error {
	return c.Terminate(ctx)
}

func WithInitialDatabase(user string, password string, dbName string) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.Env["POSTGRES_USER"] = user
		req.Env["POSTGRES_PASSWORD"] = password
		req.Env["POSTGRES_DB"] = dbName
	}
}
