package natstrace

import (
	"context"

	"github.com/nats-io/nats.go"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"

	"github.com/kabelsea-sandbox/slice/pkg/natsclient"
)

// Tracing
func Tracing(ctx context.Context, message *natsclient.Message, publish natsclient.PublishFunc) error {
	span := trace.FromContext(ctx)
	message.Data = append(propagation.Binary(span.SpanContext()), message.Data...)
	return publish(ctx, message)
}

// TraceFromMessage
func TraceFromMessage(message *nats.Msg) trace.SpanContext {
	if sc, ok := propagation.FromBinary(message.Data); ok {
		message.Data = message.Data[29:]
		return sc
	}
	return trace.SpanContext{}
}
