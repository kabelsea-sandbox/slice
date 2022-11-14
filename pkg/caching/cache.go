package caching

import (
	"context"
	"fmt"

	"github.com/eko/gocache/v3/cache"
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
)

//go:generate mockgen --package=cachingmock -destination=mocks/mock_cache.go . Cache

// Cache interface
type Cache interface {
	// Get
	Get(ctx context.Context, key any, returnObj any) (bool, error)

	// Set
	Set(ctx context.Context, key, object any) error
}

// Cache interface implementaton with MessagePack serialization
type _cache struct {
	cache cache.CacheInterface[any]
}

// New cache constructor
func New(cache cache.CacheInterface[any]) Cache {
	return &_cache{
		cache: cache,
	}
}

// Get implements Cache interface
func (c *_cache) Get(ctx context.Context, key any, returnObj any) (bool, error) {
	result, err := c.cache.Get(ctx, key)
	if err != nil {
		return false, err
	}

	if result == nil {
		return false, nil
	}

	switch v := result.(type) {
	case []byte:
		err = msgpack.Unmarshal(v, returnObj)
	case string:
		err = msgpack.Unmarshal([]byte(v), returnObj)
	default:
		err = fmt.Errorf("unexpected type, %v", v)
	}

	if err != nil {
		return false, errors.Wrap(err, "msgpack unmarshal failed")
	}
	return true, nil
}

// Set implements Cache interface
func (c *_cache) Set(ctx context.Context, key, object any) error {
	bytes, err := msgpack.Marshal(object)
	if err != nil {
		return errors.Wrap(err, "msgpack marshal failed")
	}
	return c.cache.Set(ctx, key, bytes)
}
