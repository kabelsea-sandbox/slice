package httpserver

import (
	"net/http"

	"github.com/kabelsea-sandbox/slice"
)

// Bundle
type Bundle struct {
	Addr string `envconfig:"addr" default:":8081"`
}

func (b *Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(func() *http.Server {
		return &http.Server{Addr: b.Addr}
	})
}
