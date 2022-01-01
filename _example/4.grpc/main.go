package main

import (
	"github.com/kabelsea-sanbox/slice/pkg/di"
	"google.golang.org/grpc"

	"github.com/kabelsea-sanbox/slice"
	"github.com/kabelsea-sanbox/slice/bundle/envconfig"
	grpcbundle "github.com/kabelsea-sanbox/slice/bundle/grpc"
	"github.com/kabelsea-sanbox/slice/bundle/monitoring"
	"github.com/kabelsea-sanbox/slice/bundle/zap"
)

// UserService
type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) RegisterGRPCServer(_ *grpc.Server) {

}

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