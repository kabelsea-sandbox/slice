package main

import (
	"context"

	"github.com/kabelsea-sanbox/slice/pkg/di"
	"github.com/nats-io/nats.go"

	"github.com/kabelsea-sanbox/slice"
	"github.com/kabelsea-sanbox/slice/bundle/monitoring"
	natsbundle "github.com/kabelsea-sanbox/slice/bundle/nats"
	"github.com/kabelsea-sanbox/slice/bundle/zap"
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
