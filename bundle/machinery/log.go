package machinery

import (
	"github.com/RichardKnop/logging"
	"go.uber.org/zap"
)

type LogAdapter interface {
	logging.LoggerInterface
}

type logAdapter struct {
	*zap.SugaredLogger
}

func NewLogAdapter(logger *zap.Logger) LogAdapter {
	return &logAdapter{
		logger.With(zap.Namespace("machinery")).Sugar(),
	}
}

func (l *logAdapter) Print(args ...any) {
	l.SugaredLogger.Info(args...)
}

func (l *logAdapter) Printf(s string, args ...any) {
	l.SugaredLogger.Infof(s, args...)
}

func (l *logAdapter) Println(args ...any) {
	l.SugaredLogger.Info(args...)
}
