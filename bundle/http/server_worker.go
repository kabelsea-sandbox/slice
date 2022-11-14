package http

import (
	"context"
	"net/http"

	"slice"
)

// ServerWorker control lifecycle for HTTP server.
type ServerWorker struct {
	logger slice.Logger
	server *http.Server
}

// Run runs worker.
func (p *ServerWorker) Run(context.Context) (err error) {
	p.logger.Debugf("http", "Starting http server")
	defer p.logger.Debugf("http", "Stopping http server")
	if err = p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop stops worker.
func (p *ServerWorker) Stop(err error) {
	_ = p.server.Close()
}
