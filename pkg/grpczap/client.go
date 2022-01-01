package grpczap

import (
	"context"
	"fmt"
	"path"
	"time"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	ctxzap "github.com/kabelsea-sanbox/slice/pkg/zaplog/ctx"
)

// ClientLogging is a middleware for client logging.
func ClientLogging() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		startTime := time.Now()

		service := path.Dir(method)[1:]
		callMethod := path.Base(method)

		logger := ctxzap.Extract(ctx).With(
			zap.String("grpc.system", "client"),
			zap.String("grpc.service", service),
			zap.String("grpc.method", callMethod),
		)

		logger.Debug("GRPC: Start call")

		err := invoker(ctx, method, req, reply, cc, opts...)

		code := grpc_logging.DefaultErrorToCode(err)
		level := CodeToLevel(code)

		logger.Check(level, fmt.Sprintf("GRPC: Finished call with status %s", code)).
			Write(
				zap.Error(err),
				zap.String("grpc.code", code.String()),
				zap.Duration("duration", time.Since(startTime)),
			)
		return err
	}
}
