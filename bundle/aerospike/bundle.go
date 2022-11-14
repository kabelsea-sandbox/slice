package aerospike

import (
	"github.com/aerospike/aerospike-client-go"

	"slice"
)

type Bundle struct {
	Host string `envconfig:"AEROSPIKE_HOST" required:"true"`
	Port int    `envconfig:"AEROSPIKE_PORT" required:"true" default:"3000"`
}

// Build implements Bundle.
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(b.NewClient)
}

// NewClient creates redis client.
func (b *Bundle) NewClient(logger slice.Logger) (*aerospike.Client, error) {
	logger.Infof("aerospike", "Create aerospike connection")

	return aerospike.NewClientWithPolicy(nil, b.Host, b.Port)
}
