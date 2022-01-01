package natsclient

import (
	"context"

	"github.com/nats-io/nats.go"
)

// Client is a nats client.
type Client struct {
	conn        *nats.Conn
	publishFunc PublishFunc
}

// New creates new nats client.
func New(conn *nats.Conn, options ...Option) *Client {
	client := &Client{
		conn: conn,
		publishFunc: func(ctx context.Context, message *Message) (err error) {
			return conn.Publish(message.Subject, message.Data)
		},
	}
	for _, opt := range options {
		opt.apply(client)
	}
	return client
}

// Publish publishes data to nats with subject.
func (c *Client) Publish(ctx context.Context, subject string, data []byte) (err error) {
	message := &Message{
		Subject: subject,
		Data:    data,
	}
	return c.publishFunc(ctx, message)
}

func (c *Client) publish(ctx context.Context, message *Message) (err error) {
	return c.conn.Publish(message.Subject, message.Data)
}

// Message is a nats message.
type Message struct {
	Subject string
	Data    []byte
}

// Interceptor is a nats interceptor.
type Interceptor func(ctx context.Context, message *Message, publish PublishFunc) error

// PublishFunc is a publish func.
type PublishFunc func(ctx context.Context, message *Message) (err error)
