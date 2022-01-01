package nats

import (
	"github.com/nats-io/nats.go"
)

func NewSubscriptionFactory(wrapperFactory *MessageWrapperFactory, clientID string, group string) *SubscriptionFactory {
	return &SubscriptionFactory{
		wrapperFactory: wrapperFactory,
		worker:         group,
		clientID:       clientID,
	}
}

type SubscriptionFactory struct {
	wrapperFactory *MessageWrapperFactory

	worker   string
	clientID string
}

// CreateSubscription
func (f *SubscriptionFactory) CreateSubscription(handler MessageHandler) *Subscription {
	wrapper := f.wrapperFactory.Wrap(handler)
	return NewSubscription(wrapper, f.worker)
}

type Subscription struct {
	handler      *MessageHandlerWrapper
	subscription *nats.Subscription

	worker string
}

func NewSubscription(handler *MessageHandlerWrapper, worker string) *Subscription {
	return &Subscription{
		handler: handler,
		worker:  worker,
	}
}

// Subscribe
func (s *Subscription) Subscribe(conn *nats.Conn) (err error) {
	s.subscription, err = conn.QueueSubscribe(
		s.handler.Subject(),
		s.worker,
		s.handler.Handle,
	)
	if err != nil {
		return err
	}
	return nil
}

// Unsubscribe
func (s *Subscription) Unsubscribe() (err error) {
	return s.subscription.Unsubscribe()
}
