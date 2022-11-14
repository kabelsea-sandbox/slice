package machinery

import (
	"context"
	"slice/pkg/run"

	machinery "github.com/RichardKnop/machinery/v2"
)

// Worker adapter for slice run.Worker interface
type Worker interface {
	run.Worker
}

type worker struct {
	worker      *machinery.Worker
	consumerTag string
	stop        chan struct{}
	done        chan struct{}
}

func NewWorker(
	server *machinery.Server,
	tasks TaskPool,
	consumerTag string,
	errHandler ErrorHandler,
	preHandler PreTaskHandler,
	postHandler PostTaskHandler,
) (Worker, error) {
	w := server.NewWorker(consumerTag, 0)

	w.SetErrorHandler(errHandler)
	// w.SetPreTaskHandler(preHandler)
	// w.SetPostTaskHandler(postHandler)

	for name, fn := range tasks.(*taskPool).tasks {
		if err := server.RegisterTask(name, fn); err != nil {
			return nil, err
		}
	}

	return &worker{
		worker:      w,
		consumerTag: consumerTag,
		stop:        make(chan struct{}, 1),
		done:        make(chan struct{}, 1),
	}, nil
}

// Run implements Worker interface
func (w *worker) Run(ctx context.Context) error {
	go func() {
		if err := w.worker.Launch(); err != nil {
			panic(err)
		}
	}()

	// wait until stop
	<-w.stop

	w.worker.Quit()

	close(w.done)
	return nil
}

// Stop implements Worker interface
func (w *worker) Stop(err error) {
	close(w.stop)

	// wait until done
	<-w.done
}
