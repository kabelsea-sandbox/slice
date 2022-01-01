package ctxzap

import (
	"context"

	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

// NewContext creates new context with zap logger.
func NewContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxContextKey{}, logger)
}

// SetupContext setup zap logger with setupFunc.
func SetupContext(setupFunc func(key interface{}, value interface{}), logger *zap.Logger) {
	setupFunc(ctxContextKey{}, logger)
}

// Extract extracts zap logger from context.
func Extract(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(ctxContextKey{}).(*zap.Logger)
	if !ok || logger == nil {
		return zap.NewNop()
	}

	return logger
}

// TraceID adds trace id to logger.
func TraceID(ctx context.Context) context.Context {
	traceID := trace.FromContext(ctx).SpanContext().TraceID.String()
	return context.WithValue(ctx, ctxContextKey{}, Extract(ctx).With(zap.String("trace_id", traceID)))
}

type ctxContextKey struct{}
