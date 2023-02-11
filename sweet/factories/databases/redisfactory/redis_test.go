package redisfactory_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/barry-hennessy/test/sweet"
	"github.com/barry-hennessy/test/sweet/factories/databases/redisfactory"
)

func TestNewContainerClientFactory(t *testing.T) {
	ctx := context.Background()
	redisImages := []string{
		"redis/redis-stack-server:latest",
		"redis:6",
		"redis:7",
		"redis:5",
	}

	for _, image := range redisImages {
		t.Run(fmt.Sprintf("test/%s", image), func(t *testing.T) {
			f := redisfactory.NewContainerClientFactory(ctx, image, &redis.Options{})
			sweet.Run(t, "can create and test against basic commands", f, func(t *testing.T, r *redis.Client) {
				testKey := "test"
				testVal := "val"
				r.LPush(ctx, testKey, testVal)
				val, err := r.LPop(ctx, testKey).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if val != testVal {
					t.Errorf("value was not set correctly. Expected: %q\tGot %q", testVal, val)
					return
				}
			})

			sweet.Run(t, "can publish and subscribe", f, func(t *testing.T, r *redis.Client) {
				testChan := "channel"
				testMessage := "message"

				wg := sync.WaitGroup{}
				subClient := redis.NewClient(r.Options())
				ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()

				sub := subClient.Subscribe(ctx, testChan)
				for {
					err := sub.Ping(ctx)
					if err != nil {
						t.Logf("error pinging connection: %s", err)
						continue
					}

					break
				}

				wg.Add(1)
				go func() {
					defer wg.Done()

					select {
					case received := <-sub.Channel():
						if received.Payload != testMessage {
							t.Errorf("received message was not as expected. Expected: %q\tGot %q", testMessage, received.String())
							return
						}
					case <-ctx.Done():
						t.Error(ctx.Err())
						return
					}
				}()

				published := false
				for !published {
					err := r.Publish(ctx, testChan, testMessage).Err()
					if err != nil {
						t.Logf("error publishing: %s", err)
						continue
					}
					published = true
				}

				wg.Wait()
			})
		})
	}

}

func TestNewClientFactory(t *testing.T) {
	ctx := context.Background()
	cf := redisfactory.NewClientFactory(ctx, &redis.Options{})

	t.Run("closes it's connection once done", func(t *testing.T) {
		var oldRedisClient *redis.Client
		sweet.Run(t, "do something with redis...", cf, func(t *testing.T, r *redis.Client) {
			oldRedisClient = r
		})

		_, err := oldRedisClient.Get(ctx, "key").Result()
		if !errors.Is(err, redis.ErrClosed) {
			t.Errorf("redis client should have been closed %s", err)
		}
	})
}
