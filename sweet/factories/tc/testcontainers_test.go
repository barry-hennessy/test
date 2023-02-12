package tc_test

import (
	"context"
	"testing"

	"github.com/barry-hennessy/test/sweet"
	"github.com/barry-hennessy/test/sweet/factories/tc"
	"github.com/barry-hennessy/test/sweet/factories/tc/cockroachdb"
	"github.com/barry-hennessy/test/sweet/factories/tc/redis"
	"github.com/testcontainers/testcontainers-go"
)

func TestNewFactory(t *testing.T) {
	ctx := context.Background()
	t.Run("containers start", func(t *testing.T) {
		factories := map[string]tc.Container{
			"redis:6":    redis.NewRedisContainer("redis:6"),
			"redis:7":    redis.NewRedisContainer("redis:6"),
			"crdb: 21.1": cockroachdb.NewCockroachDBContainer("cockroachdb/cockroach:latest-v21.1"),
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
