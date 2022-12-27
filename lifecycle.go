package slice

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/kabelsea-sandbox/slice/pkg/di"
)

// Initialization step of application lifecycle. It collects user dependency injection
// options and creates container with them. Incorrect dependency injection option will cause error.
func initialization(diopts ...di.Option) (*di.Container, error) {
	// create container and validate user providers
	container, err := di.New(diopts...)
	if err != nil {
		return nil, err
	}
	return container, nil
}

type envProcessFn func(prefix string, spec interface{}) error

// Configures bundle values. Now we use envconfig for this.
// It iterates over all bundles and process it with envconfig library.
func configureBundles(process envProcessFn, prefix string, bundles ...bundle) error {
	for _, bundle := range bundles {
		if err := process(prefix, bundle.Bundle); err != nil {
			return errors.Wrapf(err, "%s bundle configure failed", bundle.name)
		}
	}
	return nil
}

// Builds bundle dependencies. It iterates over all bundles and call function Build() to provide
// bundle dependencies.
func buildBundles(container *di.Container, bundles ...bundle) error {
	for _, bundle := range bundles {
		builder := newBundleContainerBuilder(container)
		bundle.Build(builder)
		if err := builder.Error(); err != nil {
			return errors.Wrapf(err, "%s: build failed", bundle.name)
		}
	}
	return nil
}

// Boots bundles.
func boot(timeout time.Duration, container *di.Container, bundles ...bundle) (shutdowns shutdowns, _ error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for _, bundle := range bundles {
		if err := ctx.Err(); err != nil {
			return shutdowns, errors.Wrapf(err, "boot failed")
		}
		if boot, ok := bundle.Bundle.(BootShutdown); ok {
			// boot bundle
			if err := boot.Boot(ctx, container); err != nil {
				return shutdowns, errors.Wrapf(err, "%s bundle boot failed", bundle.name)
			}
			// append successfully booted bundle shutdown
			shutdowns = append(shutdowns, bundleShutdown{
				name:     bundle.name,
				shutdown: boot.Shutdown,
			})
		}
	}
	return shutdowns, nil
}

// drun resolves dispatcher and run its.
func drun(container *di.Container) error {
	// resolve dispatcher
	var dispatcher Dispatcher
	if err := container.Resolve(&dispatcher); err != nil {
		return errors.Wrap(err, "resolve dispatcher failed")
	}
	// dispatcher run
	// todo: use context for stop
	if err := dispatcher.Run(context.Background()); err != nil {
		return err
	}
	return nil
}

// reverseShutdown shutdowns in reverse order.
func reverseShutdown(timeout time.Duration, container *di.Container, shutdowns shutdowns) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// shutdown bundles in reverse order
	var errs errShutdown
	for i := len(shutdowns) - 1; i >= 0; i-- {
		// bundle shutdown
		bs := shutdowns[i]
		if err := ctx.Err(); err != nil {
			return errors.Wrapf(err, "shutdown failed")
		}
		if err := bs.shutdown(ctx, container); err != nil {
			errs = append(errs, errors.Wrapf(err, "%s", bs.name))
		}
	}
	if len(errs) != 0 {
		return errors.Wrap(errs, "shutdown failed")
	}
	return nil
}

type bundleShutdown struct {
	name     string
	shutdown func(ctx context.Context, container Container) error
}

type shutdowns []bundleShutdown

type errShutdown []error

func (e errShutdown) Error() string {
	var s []string
	for _, err := range e {
		s = append(s, err.Error())
	}
	return strings.Join(s, "; ")
}
