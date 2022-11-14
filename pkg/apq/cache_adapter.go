package apq

import (
	"context"

	"go.uber.org/zap"

	"github.com/kabelsea-sandbox/slice/pkg/caching"
)

//go:generate mockgen --package=apqmock -destination=mocks/mock_cache_adapter.go . CacheAdapter

// CacheAdapter interface for APQ
type CacheAdapter interface {
	Add(ctx context.Context, key string, value any)
	Get(ctx context.Context, key string) (any, bool)
}

// cacheAdapter interface implementation
type cacheAdapter struct {
	logger *zap.Logger
	cache  caching.Cache
}

func NewCacheAdapter(logger *zap.Logger, cache CacheAPQ) CacheAdapter {
	return &cacheAdapter{
		logger: logger.With(zap.Namespace("apq")),
		cache:  cache,
	}
}

// Add implement CacheAdapter interface
func (c *cacheAdapter) Add(ctx context.Context, key string, value any) {
	logger := c.logger.With(
		zap.String("key", key),
		zap.String("value", key),
	)

	if err := c.cache.Set(ctx, key, value); err != nil {
		logger.Error("set cache failed",
			zap.Error(err),
		)
	}
}

// Add implement CacheAdapter interface
func (c *cacheAdapter) Get(ctx context.Context, key string) (any, bool) {
	var result string

	logger := c.logger.With(
		zap.String("key", key),
	)

	ok, err := c.cache.Get(ctx, key, &result)
	if err != nil {
		logger.Error("get cache failed",
			zap.Error(err),
		)
		return "", false
	}
	return result, ok
}
