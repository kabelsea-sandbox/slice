package goredis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"slice"
)

// Bundle is a redis bundle.
type Bundle struct {
	Addr     string `envconfig:"REDIS_ADDR" required:"true"`
	Database int    `envconfig:"REDIS_DATABASE" required:"true"`
	Password string `envconfig:"REDIS_PASSWORD"`
}

// Build implements Bundle.
func (b Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(b.NewClient)
}

// NewClient creates redis client.
func (b Bundle) NewClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     b.Addr,
		Password: b.Password,
		DB:       b.Database,
	})
	if err := client.Ping(context.TODO()).Err(); err != nil {
		return nil, errors.Wrap(err, "ping redis")
	}
	return client, nil
}
