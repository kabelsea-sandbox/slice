package nats

import (
	"context"
	"os"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"slice"
	"slice/pkg/di"
	"slice/pkg/natsclient"
	"slice/pkg/natstrace"
	"slice/pkg/natszap"
	"slice/pkg/run"
)

// MessageHandler
type MessageHandler interface {
	Subject() string
	Handle(ctx context.Context, message *nats.Msg) (err error)
}

// Bundle integrates NATS.
type Bundle struct {
	Hosts         string `envconfig:"nats_hosts" default:"nats://nats:4222"`
	ClusterID     string `envconfig:"nats_cluster_id"`
	ClientID      string `envconfig:"nats_client_id"`
	WorkerGroup   string `envconfig:"nats_worker_group"`
	MaxReconnects int    `envconfig:"nats_max_reconnects" default:"5"`
	ReconnectWait int    `envconfig:"nats_reconnect_wait" default:"2"`
}

// Build implements Bundle interface.
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	if b.ClientID == "" {
		b.ClientID = os.Getenv("HOSTNAME")
	}
	// nats
	builder.Provide(b.NewConnection)
	builder.Provide(b.NewClient)
	builder.Provide(NewMessageWrapperFactory)
	builder.Provide(b.NewSubscriptionFactory)
	builder.Provide(NewSubscriptionWorker, di.As(new(run.Worker)))
}

// Boot implements BootShutdown interface.
func (b *Bundle) Boot(ctx context.Context, interactor slice.Container) (err error) {
	var handlers []MessageHandler
	if interactor.Has(&handlers) {
		if err = interactor.Invoke(b.RegisterHandlers); err != nil {
			return errors.Wrap(err, "register nats message handlers")
		}
	}
	return nil
}

// Shutdown implements Bundle interface.
func (b *Bundle) Shutdown(_ context.Context, _ slice.Container) (err error) {
	return nil
}

// NewConnection creates NATS connection.
func (b *Bundle) NewConnection(logger *zap.Logger) (*nats.Conn, error) { // todo: remove *zap.Logger
	defer logger.Sugar().Info("[NATS] ", "Create nats connection")

	return nats.Connect(
		b.Hosts,
		nats.MaxReconnects(b.MaxReconnects),
		nats.ReconnectWait(time.Duration(b.ReconnectWait)*time.Second),
		nats.DisconnectHandler(func(nc *nats.Conn) {
			if err := nc.LastError(); err != nil {
				logger.Info("NATS disconnected")
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Info("NATS reconnected")
		}),
		nats.DiscoveredServersHandler(func(nc *nats.Conn) {
			logger.Info("NATS connection discovered", zap.String("nats.connected_url", nc.ConnectedUrl()))
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			if err := nc.LastError(); err != nil {
				logger.Error("NATS connection closed", zap.Error(err))

				// Exit
				_ = syscall.Kill(os.Getpid(), syscall.SIGINT)
			} else {
				logger.Info("NATS connection closed")
			}
		}),
	)
}

// NewClient creates client.
func (b *Bundle) NewClient(logger slice.Logger, conn *nats.Conn) *natsclient.Client {
	return natsclient.New(conn, natsclient.ChainInterceptor(
		natstrace.Tracing,
		natszap.Logging, // todo: remove zap integration
	))
}

// NewSubscriptionFactory creates subscription factory.
func (b *Bundle) NewSubscriptionFactory(wrapperFactory *MessageWrapperFactory) *SubscriptionFactory {
	return NewSubscriptionFactory(wrapperFactory, b.ClientID, b.WorkerGroup)
}

// RegisterHandlers registers handler.
func (b *Bundle) RegisterHandlers(worker *SubscriptionWorker, handlers []MessageHandler) {
	for _, handler := range handlers {
		worker.AddHandler(handler)
	}
}
