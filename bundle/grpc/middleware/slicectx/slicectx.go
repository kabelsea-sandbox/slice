package slicectx

import (
	"context"

	"google.golang.org/grpc"

	"github.com/kabelsea-sanbox/slice"
)

// UnaryServerInterceptor return unary server interceptor that join application context and request context.
func UnaryServerInterceptor(appctx *slice.Context) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		return handler(appctx.Join(ctx), req)
	}
}

// UnaryServerInterceptor return unary client interceptor that join application context and request context.
func UnaryClientInterceptor(appctx *slice.Context) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return invoker(appctx.Join(ctx), method, req, reply, cc, opts...)
	}
}
