package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/kabelsea-sanbox/slice"
)

// ServerWorker is a grpc server worker.
type ServerWorker struct {
	logger slice.Logger
	port   string
	server *grpc.Server
}

// Run runs grpc server worker.
func (d *ServerWorker) Run(context.Context) error {
	d.logger.Debugf("grpc-bundle", "Starting grpc server")
	defer d.logger.Debugf("grpc-bundle", "Stopping grpc server")
	l, err := net.Listen("tcp", net.JoinHostPort("", d.port))
	if err != nil {
		return err
	}
	return d.server.Serve(l)
}

// Stop stops grpc server worker.
func (d *ServerWorker) Stop(err error) {
	d.server.Stop()
}
