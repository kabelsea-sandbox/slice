package slice_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kabelsea-sandbox/slice"
)

type contextKey string

var (
	ContextKey = contextKey("ctx-key")
)

func TestContext(t *testing.T) {
	b := slice.NewContext()
	b.Set("key", "value")
	ctx := context.WithValue(context.Background(), ContextKey, "ctx-value") // TODO
	joined := b.Join(ctx)
	require.Equal(t, joined.Value("key"), "value")
	require.Equal(t, joined.Value(ContextKey), "ctx-value")
}
