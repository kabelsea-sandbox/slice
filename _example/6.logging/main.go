package main

import (
	"net"

	"go.uber.org/zap"

	"github.com/kabelsea-sandbox/slice"
	zapbundle "github.com/kabelsea-sandbox/slice/bundle/zap"
)

type foo struct {
	One string
	Two string
}

// Run
func Run(logger *zap.Logger) error {
	logger.Debug("test", zap.Any("foo", foo{One: "one", Two: "two"}))
	logger.Info("test", zap.Any("foo", foo{One: "one", Two: "two"}))
	logger.Warn("test", zap.Any("foo", foo{One: "one", Two: "two"}))
	logger.Error("test", zap.Any("foo", foo{One: "one", Two: "two"}))
	logger.Info("info",
		zap.String("user_email", "dev@example.com"),
		zap.Any("server", net.TCPAddr{}),
		zap.Any("foo", foo{One: "one", Two: "two"}),
		zap.Any("foo2", foo{One: "one", Two: "two"}),
		zap.Any("foo3", foo{One: "one", Two: "two"}),
		zap.Any("foo4", foo{One: "one", Two: "two"}),
		zap.Any("foo5", foo{One: "one", Two: "two"}),
	)
	return nil
}

func main() {
	slice.Run(
		slice.SetName("logging-example"),
		slice.SetDispatcher(Run),
		slice.RegisterBundles(
			&zapbundle.Bundle{},
		),
	)
}
