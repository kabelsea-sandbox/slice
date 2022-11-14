package slice

import (
	"time"

	"github.com/kabelsea-sandbox/slice/pkg/di"
	"github.com/kabelsea-sandbox/slice/pkg/run"
)

// Option is a application option. All diopts starts with Set prefix.
type Option interface {
	apply(options *Application)
}

// SetName sets application name.
func SetName(name string) Option {
	return sliceOption(func(app *Application) {
		app.Name = name
	})
}

// RegisterBundles is a option that registers bundles. Bundle is a reusable part of functionality, like transport,
// monitoring, logging. Bundles registers their functionality in order that it presented in this function.
func RegisterBundles(bundles ...Bundle) Option {
	return sliceOption(func(app *Application) {
		app.Bundles = append(app.Bundles, bundles...)
	})
}

// ConfigureContainer apply dependency injection container diopts.
func ConfigureContainer(options ...di.Option) Option {
	return sliceOption(func(app *Application) {
		app.di = append(app.di, options...)
	})
}

// SetStartTimeout sets start timeout.
func SetStartTimeout(timeout time.Duration) Option {
	return sliceOption(func(app *Application) {
		app.StartTimeout = timeout
	})
}

// SetStopTimeout sets stop timeout.
func SetStopTimeout(timeout time.Duration) Option {
	return sliceOption(func(app *Application) {
		app.StopTimeout = timeout
	})
}

// SetDispatcher is a option that set main invocation. It must be a function with known dependencies as parameters
// and optional error as result.
func SetDispatcher(dispatch interface{}) Option {
	return sliceOption(func(app *Application) {
		app.dispatchFunc = dispatch
	})
}

// UseWorkerDispatcher is a option that set main dispatcher and provide control flow.
func UseWorkerDispatcher() Option {
	return sliceOption(func(app *Application) {
		app.di = append(app.di,
			di.Provide(run.NewDispatcher, di.As(new(Dispatcher))),
			di.Provide(run.NewShutdowner, di.As(new(run.Worker))),
		)
	})
}

type sliceOption func(options *Application)

func (f sliceOption) apply(options *Application) { f(options) }
