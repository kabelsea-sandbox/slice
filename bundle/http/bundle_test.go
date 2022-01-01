package http

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/kabelsea-sanbox/slice/slicetest"
)

var (
	handler      http.Handler
	mux          *chi.Mux
	server       *http.Server
	serverWorker *ServerWorker
)

func TestBundle_Build(t *testing.T) {
	t.Run("bundle provides server worker with correct configuration", func(t *testing.T) {
		bundle := &Bundle{
			Port:         "1000",
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		}
		builder := slicetest.NewTestContainer(t)
		bundle.Build(builder)
		builder.ShouldNoError()
		require.True(t, builder.Has(&handler))
		require.True(t, builder.Has(&mux))
		require.True(t, builder.Has(&server))
		require.True(t, builder.Has(&serverWorker))
	})
}
