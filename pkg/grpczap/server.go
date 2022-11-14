package grpczap

import (
	"context"
	"fmt"
	"path"
	"time"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ctxzap "github.com/kabelsea-sandbox/slice/pkg/zaplog/ctx"
)

// ServerLogging is a middleware for server logging.
func ServerLogging() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()
		service := path.Dir(info.FullMethod)[1:]
		method := path.Base(info.FullMethod)
		span := trace.FromContext(ctx)
		logger := ctxzap.Extract(ctx).With(
			zap.String("grpc.system", "server"),
			zap.String("grpc.service", service),
			zap.String("grpc.method", method),
			zap.String("trace_id", span.SpanContext().TraceID.String()),
		)
		logger.Debug("GRPC: Received")
		resp, err = handler(ctxzap.TraceID(ctx), req)
		code := grpc_logging.DefaultErrorToCode(err)
		level := CodeToLevel(code)
		logger.Check(level, fmt.Sprintf("GRPC: Processed with status %s", code)).
			Write(
				zap.Error(err),
				zap.String("grpc.code", code.String()),
				zap.Duration("duration", time.Since(start)),
			)
		return resp, err
	}
}

// Recovery returns a new unary server interceptor for panic recovery.
func Recovery() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (_ interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				ctxzap.Extract(ctx).Error("GRPC: Panic recovered", zap.Reflect("panic", r))
				err = status.Errorf(codes.Internal, "panic: %s", r)
			}
		}()
		return handler(ctx, req)
	}
}
