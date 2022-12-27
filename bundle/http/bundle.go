package http

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"

	"github.com/kabelsea-sandbox/slice/pkg/run"

	"github.com/kabelsea-sandbox/slice/pkg/di"

	"github.com/kabelsea-sandbox/slice"
)

// Controller
type Controller interface {
	// RegisterRoutes registers controller routes via chi.Router.
	RegisterRoutes(mux chi.Router)
}

// Bundle provides http functionality
type Bundle struct {
	Port         string        `envconfig:"HTTP_PORT" default:"8080"`
	ReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"5s"`
	WriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"5s"`
}

// Build provides bundle dependencies.
func (b Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(b.NewRouter, di.As(new(http.Handler)))
	builder.Provide(b.NewServer)
	builder.Provide(b.NewServerWorker, di.As(new(run.Worker)))
}

func (b Bundle) Boot(_ context.Context, container slice.Container) (err error) {
	var mux *chi.Mux
	if err := container.Resolve(&mux); err != nil {
		return err
	}

	var middlewares []func(http.Handler) http.Handler

	if container.Has(&middlewares) {
		if err := container.Resolve(&middlewares); err != nil {
			return err
		}
		mux.Use(middlewares...)
	}

	var controllers []Controller

	if container.Has(&controllers) {
		if err := container.Resolve(&controllers); err != nil {
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
	router := chi.NewRouter()

	router.Use(cors.AllowAll().Handler)

	// router.Use(cors.New(cors.Options{
	// 	AllowedOrigins: []string{"*"},
	// 	// AllowCredentials: true,
	// 	// Debug:            true,
	// }).Handler)
	return router
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
