package aerospike

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/pkg/errors"

	"github.com/kabelsea-sanbox/slice"
)

type Bundle struct {
	Host string `envconfig:"AEROSPIKE_HOST" required:"true"`
	Port int    `envconfig:"AEROSPIKE_PORT" required:"true" default:"3000"`
}

// Build implements Bundle.
func (b Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(b.NewClient)
}

// NewClient creates redis client.
func (b Bundle) NewClient() (*aerospike.Client, error) {
	client, err := aerospike.NewClientWithPolicy(nil, b.Host, b.Port)
	if err != nil {
		return nil, errors.Wrap(err, "aerospike client")
	}
	return client, nil
}
