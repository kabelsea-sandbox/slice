package natszap

import (
	"context"

	"go.uber.org/zap"

	ctxzap "github.com/kabelsea-sandbox/slice/pkg/zaplog/ctx"

	"github.com/kabelsea-sandbox/slice/pkg/natsclient"
)

// Logging
func Logging(ctx context.Context, message *natsclient.Message, publish natsclient.PublishFunc) error {
	logger := ctxzap.Extract(ctx).With(
		zap.String("nats.subject", message.Subject),
	)

	logger.Debug("NATS: Publish")
	if err := publish(ctx, message); err != nil {
		logger.Error("NATS: Publish error", zap.Error(err))
	}

	logger.Info("NATS: Published")
	return nil
}
