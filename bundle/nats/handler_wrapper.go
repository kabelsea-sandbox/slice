package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"go.opencensus.io/trace"
	"go.uber.org/zap"

	"github.com/kabelsea-sanbox/slice"
	"github.com/kabelsea-sanbox/slice/pkg/natstrace"
	ctxzap "github.com/kabelsea-sanbox/slice/pkg/zaplog/ctx"
)

type MessageWrapperFactory struct {
	ctx *slice.Context
}

// NewMessageWrapperFactory
func NewMessageWrapperFactory(ctx *slice.Context) *MessageWrapperFactory {
	return &MessageWrapperFactory{ctx: ctx}
}

// Wrap
func (f *MessageWrapperFactory) Wrap(handler MessageHandler) *MessageHandlerWrapper {
	return &MessageHandlerWrapper{
		ctx:     f.ctx,
		handler: handler,
	}
}

type MessageHandlerWrapper struct {
	ctx     *slice.Context
	handler MessageHandler
}

// Subject
func (h *MessageHandlerWrapper) Subject() string {
	return h.handler.Subject()
}

// Handle
func (h *MessageHandlerWrapper) Handle(message *nats.Msg) {
	ctx := h.ctx.Join(context.Background())
	ctx, _ = trace.StartSpanWithRemoteParent(
		ctx, fmt.Sprintf("%s received", message.Subject), natstrace.TraceFromMessage(message),
	)
	ctx = ctxzap.TraceID(ctx)
	msglog := ctxzap.Extract(ctx).With(
		zap.String("nats.subject", message.Subject),
	)
	msglog.Debug("NATS: Received")
	if err := h.handler.Handle(ctx, message); err != nil {
		msglog.Error("NATS: Handle error", zap.Error(err))
		return
	}
	msglog.Info("NATS: Processed")
}
