package slice

import (
	"context"

	"github.com/kabelsea-sanbox/slice/pkg/di"
)

//go:generate mockgen -package=slice -destination=./dispatcher_mock_test.go -source=./dispatcher.go Dispatcher

// Dispatcher runs application.
type Dispatcher interface {
	Run(ctx context.Context) error
}

type invokeDispatcher struct {
	fn        interface{}
	container *di.Container
}

func (f invokeDispatcher) Run(ctx context.Context) error {
	return f.container.Invoke(f.fn)
}
