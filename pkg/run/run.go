package run

import (
	"context"

	"github.com/oklog/run"
)

// An worker is a run actor. Each actor in the run group will be run concurrently.
type Worker interface {
	// Run runs the actor and must block it until done.
	Run(ctx context.Context) error
	// Stop gives a stop signal for actor.
	Stop(err error)
}

// Dispatcher is a concurrent services helper.
type Dispatcher struct {
	group   run.Group
	workers []Worker
}

// NewDispatcher creates new worker dispatcher.
func NewDispatcher(workers []Worker) *Dispatcher {
	return &Dispatcher{
		group:   run.Group{},
		workers: workers,
	}
}

// Run runs all concurrent workers.
func (d *Dispatcher) Run(ctx context.Context) error {
	for _, worker := range d.workers {
		w := worker
		wrappedRun := func() error { return w.Run(ctx) }
		d.group.Add(wrappedRun, worker.Stop)
	}
	return d.group.Run()
}
