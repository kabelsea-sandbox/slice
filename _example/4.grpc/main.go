package main

import (
	"github.com/kabelsea-sandbox/slice"
	"github.com/kabelsea-sandbox/slice/bundle/envconfig"
	"github.com/kabelsea-sandbox/slice/bundle/monitoring"
	"github.com/kabelsea-sandbox/slice/bundle/zap"
	"github.com/kabelsea-sandbox/slice/pkg/di"

	"google.golang.org/grpc"

	"github.com/kabelsea-sandbox/slice"
	grpcbundle "github.com/kabelsea-sandbox/slice"
)

// UserService
type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) RegisterGRPCServer(_ *grpc.Server) {}

func main() {
	slice.Run(
		slice.SetName("grpc-example"),
		slice.UseWorkerDispatcher(),
		slice.RegisterBundles(
			&envconfig.Bundle{},
			&grpcbundle.Bundle{},
			&monitoring.Bundle{
				MetricsEnabled: true,
				TraceEnabled:   true,
			},
			&zap.Bundle{},
		),
		slice.ConfigureContainer(
			di.Provide(func() *Options { return &Options{} }, di.As(new(envconfig.Options))),
			di.Provide(NewUserService, di.As(new(grpcbundle.Service))),
		),
	)
}

type Options struct {
}
