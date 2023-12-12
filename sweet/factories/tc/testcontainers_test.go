package tc_test

import (
	"context"
	"runtime"
	"testing"

	"github.com/barry-hennessy/test/sweet"
	"github.com/barry-hennessy/test/sweet/factories/tc"
	"github.com/barry-hennessy/test/sweet/factories/tc/cockroachdb"
	"github.com/barry-hennessy/test/sweet/factories/tc/mongodb"
	"github.com/barry-hennessy/test/sweet/factories/tc/nats"
	"github.com/barry-hennessy/test/sweet/factories/tc/postgres"
	"github.com/barry-hennessy/test/sweet/factories/tc/redis"
	"github.com/testcontainers/testcontainers-go"
)

func TestNewFactory(t *testing.T) {
	if runtime.GOOS == "windows" {
		// See https://github.com/barry-hennessy/test/issues/11
		t.Skip()
		return
	}

	ctx := context.Background()
	t.Run("containers start", func(t *testing.T) {
		factories := map[string]tc.Container{
			"redis:6":     redis.NewRedisContainer("redis:6"),
			"redis:7":     redis.NewRedisContainer("redis:6"),
			"crdb:21.1":   cockroachdb.NewCockroachDBContainer("cockroachdb/cockroach:latest-v21.1"),
			"mongo:6":     mongodb.NewMongoDBContainer("mongo:6"),
			"postgres:11": postgres.NewPostgresContainer("postgres:11"),
			"postgres:12": postgres.NewPostgresContainer("postgres:12"),
			"postgres:13": postgres.NewPostgresContainer("postgres:13"),
			"postgres:14": postgres.NewPostgresContainer("postgres:14"),
			"postgres:15": postgres.NewPostgresContainer("postgres:15"),
			"nats:2":      nats.NewNatsContainer("nats:2"),
		}

		for name, container := range factories {
			factory := tc.NewFactory(ctx, container)
			sweet.Run(t, name, factory, func(t *testing.T, c testcontainers.Container) {
				state, err := c.State(ctx)
				if err != nil {
					t.Errorf("Could not get container state: %q", err)
					return
				}

				if !state.Running {
					t.Errorf("Container is not running")
				}

				if state.Error != "" {
					t.Errorf("Container is in an error state: %q", state.Error)
				}
			})
		}
	})
}
