package http

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/kabelsea-sanbox/slice"
	"github.com/kabelsea-sanbox/slice/pkg/di"
	"github.com/kabelsea-sanbox/slice/pkg/run"
)

// Controller
type Controller interface {
	// RegisterRoutes registers controller routes via chi.Router.
	RegisterRoutes(mux chi.Router)
}

// Bundle provides http functionality
type Bundle struct {
	Port         string        `envconfig:"HTTP_PORT" default:"8080"`
	ReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"2s"`
	WriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"2s"`
}

// Build provides bundle dependencies.
func (b Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(b.NewRouter, di.As(new(http.Handler)))
	builder.Provide(b.NewServer)
	builder.Provide(b.NewServerWorker, di.As(new(run.Worker)))
}

func (b Bundle) Boot(_ context.Context, container slice.Container) (err error) {
	var controllers []Controller
	if container.Has(&controllers) {
		if err := container.Resolve(&controllers); err != nil {
			return err
		}
		var mux *chi.Mux
		if err := container.Resolve(&mux); err != nil {
			return err
		}
		for _, ctrl := range controllers {
			ctrl.RegisterRoutes(mux)
		}
	}
	return nil
}

func (b Bundle) Shutdown(ctx context.Context, container slice.Container) (err error) {
	return nil
}

// NewRouter creates http router.
func (b Bundle) NewRouter() *chi.Mux {
	return chi.NewRouter()
}

// NewServer creates http server
func (b Bundle) NewServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         net.JoinHostPort("", b.Port),
		Handler:      handler,
		ReadTimeout:  b.ReadTimeout,
		WriteTimeout: b.WriteTimeout,
	}
}

// NewServerWorker creates server worker.
func (b Bundle) NewServerWorker(logger slice.Logger, server *http.Server) *ServerWorker {
	return &ServerWorker{
		logger: logger,
		server: server,
	}
}
