package graphql

import (
	"context"
	"net/http"
	"time"

	"github.com/kabelsea-sandbox/slice"
	httpbundle "github.com/kabelsea-sandbox/slice/bundle/http"
	"github.com/kabelsea-sandbox/slice/pkg/apq"
	"github.com/kabelsea-sandbox/slice/pkg/di"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/apollotracing"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type Bundle struct {
	Playground    bool `envconfig:"GRAPHQL_PLAYGROUND_ENABLED" default:"False"`
	Tracing       bool `envconfig:"GRAPHQL_TRACING_ENABLED" default:"False"`
	Intorspection bool `envconfig:"GRAPHQL_INTROSPECTION_ENABLED" default:"False"`

	Websocket             bool
	WebsocketPingInterval time.Duration `envconfig:"GRAPHQL_WEBSOCKET_PING_INTERVAL" default:"5s"`

	APQ bool
}

// DependOn
func (b *Bundle) DependOn() []slice.Bundle {
	return []slice.Bundle{
		&httpbundle.Bundle{},
	}
}

// Build provides exporters and worker to di container.
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(b.NewServer)

	builder.Provide(NewQueryController, di.As(new(httpbundle.Controller)))

	if b.Playground {
		builder.Provide(NewPlaygroundController, di.As(new(httpbundle.Controller)))
	}
}

// Boot bundle
func (b *Bundle) Boot(_ context.Context, container slice.Container) (err error) {
	var errorPresenter graphql.ErrorPresenterFunc
	if container.Has(&errorPresenter) {
		if err = container.Invoke(b.RegisterErrorPresenter); err != nil {
			return errors.Wrap(err, "graphql error presenter register failed")
		}
	}

	var websocket []*transport.Websocket
	if container.Has(&websocket) {
		if err = container.Invoke(b.RegisterWebsocketTransport); err != nil {
			return errors.Wrap(err, "graphql register websocket transport failed")
		}
	}

	return nil
}

// Shutdown implements Bundle interface.
func (b *Bundle) Shutdown(_ context.Context, _ slice.Container) (err error) {
	return nil
}

func (b *Bundle) NewServer(schema graphql.ExecutableSchema, cacheAPQ apq.CacheAdapter) *handler.Server {
	server := handler.New(schema)

	server.AddTransport(transport.Options{})
	server.AddTransport(transport.POST{})
	server.AddTransport(transport.GET{})
	server.AddTransport(transport.MultipartForm{})

	// tracing, do not use on production
	if b.Tracing {
		server.Use(apollotracing.Tracer{})
	}

	// graphql introspection, do not use on production
	if b.Intorspection {
		server.Use(extension.Introspection{})
	}

	// enabled automatic persistence query
	if b.APQ {
		server.Use(
			extension.AutomaticPersistedQuery{
				Cache: cacheAPQ,
			},
		)
	}

	// websocket transport
	if b.Websocket {
		server.AddTransport(transport.Websocket{
			Upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
			},
			PingPongInterval:      b.WebsocketPingInterval,
			KeepAlivePingInterval: b.WebsocketPingInterval,
		})
	}
	return server
}

func (b *Bundle) RegisterErrorPresenter(server *handler.Server, errorPresenter graphql.ErrorPresenterFunc) {
	server.SetErrorPresenter(errorPresenter)
}

func (b *Bundle) RegisterWebsocketTransport(server *handler.Server, transport *transport.Websocket) {
	server.AddTransport(transport)
}
