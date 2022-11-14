package apq

import "slice/pkg/caching"

//go:generate mockgen --package=apqmock -destination=mocks/mock_cache_apq.go . CacheAPQ

// CacheAPQ provide cache manager API for APQ
type CacheAPQ caching.Cache
