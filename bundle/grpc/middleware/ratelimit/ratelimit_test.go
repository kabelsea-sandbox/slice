package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDefaultLimiter_Limit(t *testing.T) {
	t.Run("limit true after 10 per second", func(t *testing.T) {
		limiter := NewDefaultLimiter(time.Second, 10)
		for i := 0; i < 10; i++ {
			require.False(t, limiter.Limit())
		}
		require.True(t, limiter.Limit())
		<-time.After(1 * time.Second)
		require.False(t, limiter.Limit())
	})
}
