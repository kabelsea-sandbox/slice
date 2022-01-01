package natsclient

import (
	"context"
)

// Option is a client option.
type Option interface {
	apply(c *Client)
}

// ChainInterceptor is chained interceptor.
func ChainInterceptor(interceptors ...Interceptor) Option {
	n := len(interceptors)
	return clientOption(func(c *Client) {
		chainer := func(currentInter Interceptor, handler PublishFunc) PublishFunc {
			return func(currentCtx context.Context, message *Message) error {
				return currentInter(currentCtx, message, handler)
			}
		}
		chainedPublishFunc := c.publish
		for i := n - 1; i >= 0; i-- {
			chainedPublishFunc = chainer(interceptors[i], chainedPublishFunc)
		}
		c.publishFunc = chainedPublishFunc
	})
}

type clientOption func(c *Client)

func (o clientOption) apply(c *Client) { o(c) }
