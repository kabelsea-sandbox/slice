package run

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Shutdowner
type Shutdowner struct {
	stop chan os.Signal
}

// NewShutdowner
func NewShutdowner() *Shutdowner {
	return &Shutdowner{
		stop: make(chan os.Signal),
	}
}

func (s *Shutdowner) Run(context.Context) error {
	signal.Notify(s.stop, syscall.SIGTERM, syscall.SIGINT)
	<-s.stop
	return nil
}

func (s *Shutdowner) Stop(err error) {
	close(s.stop)
}
