package natsjsm

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"

	natsbundle "github.com/kabelsea-sandbox/slice/bundle/nats"

	"github.com/kabelsea-sandbox/slice"
)

// Bundle integrates Nats Jetstream.
type Bundle struct {
	Timeout time.Duration `envconfig:"NATS_JETSTREAM_TIMEOUT" default:"5s"`
}

// DependOn
func (b *Bundle) DependOn() []slice.Bundle {
	return []slice.Bundle{
		&natsbundle.Bundle{},
	}
}

// Build implements Bundle interface.
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(b.NewConnection)
}

// Boot implements BootShutdown interface.
func (b *Bundle) Boot(ctx context.Context, interactor slice.Container) (err error) {
	var streams []*nats.StreamConfig

	// streams
	if interactor.Has(&streams) {
		if err = interactor.Invoke(b.RegisterStreams); err != nil {
			return errors.Wrap(err, "failed register nats stream")
		}
	}
	return nil
}

// Shutdown implements Bundle interface.
func (b *Bundle) Shutdown(_ context.Context, _ slice.Container) (err error) {
	return nil
}

// NewConnection creates Nats Jetstream connection
func (b *Bundle) NewConnection(logger slice.Logger, nc *nats.Conn) (nats.JetStreamContext, error) {
	defer logger.Infof("natjsm", "Create Jetstream connection")

	return nc.JetStream()
}

// RegisterStreams registers handler.
func (b *Bundle) RegisterStreams(worker *Worker, streams []*nats.StreamConfig) error {
	for _, stream := range streams {
		if err := worker.AddStream(stream); err != nil {
			return err
		}
	}
	return nil
}
