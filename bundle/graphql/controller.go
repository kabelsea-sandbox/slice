package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	chi "github.com/go-chi/chi/v5"
)

const (
	defaultPlaygroundPath = "/playground"
	defaultQueryPath      = "/graphql"
)

// PlaygroundController
type PlaygroundController struct{}

// NewPlaygroundController constructs controller.
func NewPlaygroundController() *PlaygroundController {
	return &PlaygroundController{}
}

// RegisterRoutes
func (c PlaygroundController) RegisterRoutes(mux chi.Router) {
	mux.Handle(defaultPlaygroundPath, playground.Handler("GraphQL Playground", defaultQueryPath))
}

// QueryController
type QueryController struct {
	server *handler.Server
}

// NewQueryController constructs controller.
func NewQueryController(server *handler.Server) *QueryController {
	return &QueryController{
		server: server,
	}
}

// RegisterRoutes
func (c *QueryController) RegisterRoutes(mux chi.Router) {
	mux.Handle(defaultQueryPath, c.server)
}
