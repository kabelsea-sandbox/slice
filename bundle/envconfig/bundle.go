package envconfig

import (
	"context"

	"github.com/kelseyhightower/envconfig"

	"slice"
)

// Options is a interface that will be loaded by envconfig bundle on boot stage and processed via envconfig.
type Options interface {
}

// Bundle is a envconfig bundle.
type Bundle struct {
	Prefix string
}

// Build nothing to inject.
func (b *Bundle) Build(_ slice.ContainerBuilder) {}

// Boot boots all envconfig options and process it via envconfig library.
func (b *Bundle) Boot(_ context.Context, container slice.Container) (err error) {
	var options []Options
	if container.Has(&options) {
		if err = container.Invoke(b.ProcessOptions); err != nil {
			return err
		}
	}
	return nil
}

// Shutdown implements Bundle interface.
func (b *Bundle) Shutdown(ctx context.Context, container slice.Container) (err error) {
	return nil
}

// ProcessOptions processes envconfig options.
func (b *Bundle) ProcessOptions(options []Options) error {
	for _, opt := range options {
		if err := envconfig.Process(b.Prefix, opt); err != nil {
			return err
		}
	}
	return nil
}
