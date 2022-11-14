package natsjsm

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Worker manage nats stream subscriptions
type Worker struct {
	logger *zap.Logger
	conn   nats.JetStreamContext
	subs   []*nats.Subscription
	stop   chan struct{}
}

// NewWorker constructor
func NewWorker(logger *zap.Logger, conn nats.JetStreamContext) *Worker {
	return &Worker{
		logger: logger,
		conn:   conn,
		stop:   make(chan struct{}),
	}
}

// AddStream create nats stream if does not exist
func (w *Worker) AddStream(stream *nats.StreamConfig) error {
	info, _ := w.conn.StreamInfo(stream.Name)

	// skip if stream exist
	if info != nil {
		w.logger.
			Sugar().
			Infof("[NATSJSM-BUNDLE] %s - %s", "Skip Nats Stream creation, exist", stream.Name)
		return nil
	}

	if _, err := w.conn.AddStream(stream); err != nil {
		return errors.Wrap(err, "nats stream add failed")
	}

	w.logger.
		Sugar().
		Infof("[NATSJSM-BUNDLE] %s - %s", "Register NATS Stream", stream.Name)

	return nil
}

func (w *Worker) AddSubscription(sub *nats.Subscription) {
	w.subs = append(w.subs, sub)
}

// Run implement Worker interface
func (w *Worker) Run(ctx context.Context) error {
	w.logger.
		Sugar().
		Info("[NATSJSM-BUNDLE]", "Starting NATS Jetstream subscription listener")

	defer w.logger.
		Sugar().
		Info("[NATSJSM-BUNDLE]", "Stopping NATS Jetstream subscription listener")

	<-w.stop

	// unsubscribe from all streams
	for _, s := range w.subs {
		if err := s.Unsubscribe(); err != nil {
			w.logger.Warn("nats stream unsubscribe failed", []zap.Field{
				zap.String("subject", s.Subject),
				zap.Error(err),
			}...)
		}
	}
	return nil
}

// Stop implement Worker interface
func (w *Worker) Stop(err error) {
	close(w.stop)
}
