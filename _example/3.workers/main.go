package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kabelsea-sandbox/slice"
	"github.com/kabelsea-sandbox/slice/bundle/zap"
	"github.com/kabelsea-sandbox/slice/pkg/di"
	"github.com/kabelsea-sandbox/slice/pkg/run"
)

// Application
type Application struct {
	WorkerCount int `envconfig:"worker_count" default:"2"`
}

func (a *Application) Build(builder slice.ContainerBuilder) {
	for i := 0; i < a.WorkerCount; i++ {
		current := i
		builder.Provide(func() *Worker {
			return &Worker{
				name: fmt.Sprintf("%d", current),
				stop: make(chan struct{}),
			}
		},
			di.WithName(fmt.Sprintf("%d", current)),
			di.As(new(run.Worker)),
		)
	}
	return
}

func main() {
	slice.Run(
		slice.SetName("workers-example"),
		slice.UseWorkerDispatcher(),
		slice.RegisterBundles(&zap.Bundle{}),
		slice.RegisterBundles(
			&Application{},
		),
	)
}

// Worker
type Worker struct {
	name string
	stop chan struct{}
}

func (w *Worker) Name() string {
	return w.name
}

func (w *Worker) Run(context.Context) error {
	for {
		select {
		case <-w.stop:
			return nil
		default:
			fmt.Println(w.name)
			time.Sleep(1 * time.Second)
		}
	}
}

func (w *Worker) Stop(err error) {
	close(w.stop)
	return
}
