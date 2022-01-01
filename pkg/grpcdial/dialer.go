package grpcdial

import (
	"github.com/sercand/kuberesolver"
	"google.golang.org/grpc"
)

// Dialer is a grpc connection dialer.
type Dialer struct {
	options []grpc.DialOption
}

// NewDialer creates new grpc dialer with default dial options.
func NewDialer(options ...grpc.DialOption) *Dialer {
	kuberesolver.RegisterInCluster()
	return &Dialer{
		options: options,
	}
}

// Dial connects with host uses default and override options.
func (f *Dialer) Dial(host string, overrides ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(host, append(f.options, overrides...)...)
}
