package grpc

import (
	"context"
	"reflect"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"slice"
	"slice/bundle/grpc/middleware/ratelimit"
	"slice/bundle/grpc/middleware/slicectx"
	"slice/bundle/monitoring"
	"slice/pkg/di"
	"slice/pkg/grpcdial"
	"slice/pkg/grpczap"
	"slice/pkg/run"
)

// Service is a grpc service interface. It will be loaded on slice boot stage and registers exists services.
type Service interface {
	// RegisterGRPCServer registers grpc server.
	RegisterGRPCServer(server *grpc.Server)
}

// Bundle is a bundle that provides configured grpc server.
type Bundle struct {
	Port                 string `envconfig:"GRPC_PORT" default:"8090"`
	MaxConcurrentStreams uint32 `envconfig:"GRPC_MAX_CONCURRENT_STREAMS" default:"250"`
	RateLimit            int    `envconfig:"GRPC_RATE_LIMIT" default:"1000"` // limit per rate period
	RateBurst            int    `envconfig:"GRPC_RATE_BURST"`                // default limit / 10
	Reflection           bool   `envconfig:"GRPC_REFLECTION" default:"True"`
}

// Build provides dialer, server and worker to di container.
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	// dialer
	builder.Provide(b.NewDialer)
	// server
	builder.Provide(newInterceptorRegistry)
	builder.Provide(b.NewServer)
	builder.Provide(b.NewServerWorker, di.As(new(run.Worker)))
	builder.Provide(b.NewDefaultRateLimiter, di.As(new(ratelimit.Limiter)))
	// integrate with monitoring bundle
	builder.Provide(b.NewMetricViews, di.As(monitoring.IMetricViews))
}

// Boot registers exists grpc services.
func (b *Bundle) Boot(_ context.Context, container slice.Container) (err error) {
	var ir *interceptorRegistry
	if err = container.Resolve(&ir); err != nil {
		return err
	}
	if err = container.Invoke(ir.SliceContext); err != nil {
		return err
	}
	if err = container.Invoke(ir.Logging); err != nil {
		return err
	}
	if err = container.Invoke(ir.Recovery); err != nil {
		return err
	}
	if err = container.Invoke(ir.RateLimit); err != nil {
		return err
	}

	var services []Service

	if container.Has(&services) {
		if err = container.Invoke(b.RegisterGRPCServices); err != nil {
			return errors.Wrap(err, "register grpc service failed")
		}
	}
	return nil
}

// Shutdown implements Bundle interface.
func (b *Bundle) Shutdown(_ context.Context, _ slice.Container) (err error) {
	return nil
}

// NewDialer creates GRPC dialer.
func (b *Bundle) NewDialer(ctx *slice.Context) *grpcdial.Dialer {
	return grpcdial.NewDialer(
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				slicectx.UnaryClientInterceptor(ctx),
				grpczap.ClientLogging(), // todo: remove grpczap link, integrate via interface.
			),
		),
		grpc.WithInsecure(),
		// grpc.WithDefaultServiceConfig(roundrobin.Name), // not working
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
	)
}

// NewServer creates GRPC server.
func (b *Bundle) NewServer(registry *interceptorRegistry) *grpc.Server {
	interceptors := registry.All()

	options := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(b.MaxConcurrentStreams),
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(interceptors...),
		),
	}

	return grpc.NewServer(options...)
}

// NewServerWorker creates server worker.
func (b *Bundle) NewServerWorker(logger slice.Logger, server *grpc.Server) *ServerWorker {
	return &ServerWorker{logger: logger, port: b.Port, server: server}
}

// NewDefaultRateLimiter creates default rate limiter.
func (b *Bundle) NewDefaultRateLimiter() *ratelimit.DefaultLimiter {
	period := time.Second / time.Duration(b.RateLimit)
	if b.RateBurst == 0 {
		b.RateBurst = b.RateLimit / 10
	}
	return ratelimit.NewDefaultLimiter(period, b.RateBurst)
}

// NewMetricViews creates metric views.
func (b *Bundle) NewMetricViews() *MetricViews {
	return NewMetricViews()
}

// RegisterGRPCServices registers GRPC services.
func (b *Bundle) RegisterGRPCServices(logger slice.Logger, server *grpc.Server, services []Service) {
	for _, service := range services {
		logger.Infof("grpc", "Register service: %s", reflect.TypeOf(service))
		service.RegisterGRPCServer(server)
	}

	// Enable reflection
	if b.Reflection && len(services) > 0 {
		logger.Infof("grpc", "Reflection enabled")
		reflection.Register(server)
	}
}
