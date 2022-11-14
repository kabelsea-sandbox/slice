package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"

	"slice"
	ctxzap "slice/pkg/zaplog/ctx"
)

func NewSubscriptionWorker(logger slice.Logger, conn *nats.Conn, factory *SubscriptionFactory) *SubscriptionWorker {
	proc := &SubscriptionWorker{
		logger:  logger,
		conn:    conn,
		factory: factory,
		done:    make(chan struct{}),
	}

	return proc
}

// SubscriptionWorker
type SubscriptionWorker struct {
	logger  slice.Logger
	conn    *nats.Conn
	factory *SubscriptionFactory
	subs    []*Subscription
	done    chan struct{}
}

func (p *SubscriptionWorker) AddHandler(handler MessageHandler) {
	p.subs = append(p.subs, p.factory.CreateSubscription(handler))
}

func (p *SubscriptionWorker) Run(ctx context.Context) error {
	p.logger.Debugf("nats", "Starting NATS subscription listener")
	defer p.logger.Debugf("nats", "Stopping NATS subscription listener")
	logger := ctxzap.Extract(ctx)
	for _, sub := range p.subs {
		logger.Debug(fmt.Sprintf("NATS Subscribe: %s", sub.handler.Subject()))
		if err := sub.Subscribe(p.conn); err != nil {
			return err
		}
	}
	<-p.done
	p.conn.Close()
	return nil
}

// Stop
func (p *SubscriptionWorker) Stop(err error) {
	close(p.done)
}
