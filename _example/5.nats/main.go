package main

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/kabelsea-sandbox/slice"
	"github.com/kabelsea-sandbox/slice/bundle/monitoring"
	natsbundle "github.com/kabelsea-sandbox/slice/bundle/nats"
	"github.com/kabelsea-sandbox/slice/bundle/zap"
	"github.com/kabelsea-sandbox/slice/pkg/di"
)

type MessageHandler struct {
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

func (m *MessageHandler) Subject() string {
	return "message-handler"
}

func (m *MessageHandler) Handle(ctx context.Context, message *nats.Msg) (err error) {
	return nil
}

func main() {
	slice.Run(
		slice.SetName("nats-example"),
		slice.UseWorkerDispatcher(),
		slice.RegisterBundles(
			&natsbundle.Bundle{},
			&monitoring.Bundle{},
			&zap.Bundle{},
		),
		slice.ConfigureContainer(
			di.Provide(NewMessageHandler, di.As(new(natsbundle.MessageHandler))),
		),
	)
}
