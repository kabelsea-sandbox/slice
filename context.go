package slice

import (
	"context"
	"sync"
	"time"
)

// Context is a slice mutable context.
type Context struct {
	parent context.Context
	values sync.Map
}

// NewContext creates new context.
func NewContext() *Context {
	return &Context{}
}

// Deadline implements context.Context interface.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	if c.parent != nil {
		return c.parent.Deadline()
	}
	return
}

// Done implements context.Context interface.
func (c *Context) Done() <-chan struct{} {
	if c.parent != nil {
		return c.parent.Done()
	}
	return nil
}

// Err implements context.Context interface.
func (c *Context) Err() error {
	if c.parent != nil {
		return c.parent.Err()
	}
	return nil
}

// Value implements context.Context interface.
func (c *Context) Value(key interface{}) interface{} {
	if c.parent != nil {
		v := c.parent.Value(key)
		if v != nil {
			return v
		}
	}
	v, _ := c.values.Load(key)
	return v
}

// Set set value to context with key.
func (c *Context) Set(key interface{}, value interface{}) {
	c.values.Store(key, value)
}

// Join merges two contexts.
func (c *Context) Join(ctx context.Context) context.Context {
	return &Context{
		parent: ctx,
		values: c.values, // nolint:govet
	}
}
