package grpc

import (
	"google.golang.org/grpc"

	"github.com/kabelsea-sandbox/slice"
	"github.com/kabelsea-sandbox/slice/bundle/grpc/middleware/ratelimit"
	"github.com/kabelsea-sandbox/slice/bundle/grpc/middleware/slicectx"
	"github.com/kabelsea-sandbox/slice/pkg/grpczap"
)

// interceptorRegistry registers and store interceptors for GRPC server.
type interceptorRegistry struct {
	server []grpc.UnaryServerInterceptor
}

// newInterceptorRegistry creates new server interceptor registry.
func newInterceptorRegistry() *interceptorRegistry {
	return &interceptorRegistry{
		server: []grpc.UnaryServerInterceptor{},
	}
}

// Register registers interceptor for GRPC server.
func (r *interceptorRegistry) Register(interceptors ...grpc.UnaryServerInterceptor) {
	r.server = append(r.server, interceptors...)
}

// All returns all registered interceptors.
func (r *interceptorRegistry) All() []grpc.UnaryServerInterceptor {
	return r.server
}

func (r *interceptorRegistry) SliceContext(ctx *slice.Context) {
	r.Register(slicectx.UnaryServerInterceptor(ctx))
}

func (r *interceptorRegistry) Logging() {
	r.Register(grpczap.ServerLogging())
}

func (r *interceptorRegistry) Recovery(env slice.Env) {
	if env.IsProduction() {
		r.Register(grpczap.Recovery())
	}
}

func (r *interceptorRegistry) RateLimit(limiter ratelimit.Limiter) {
	r.Register(ratelimit.UnaryServerInterceptor(limiter))
}
