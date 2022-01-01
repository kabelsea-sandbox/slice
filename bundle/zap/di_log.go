package zap

import (
	"fmt"

	"go.uber.org/zap"
)

// DIAdapter implements di.Logger
type DIAdapter struct {
	logger *zap.Logger
}

// NewDIAdapter create new di logger implementation.
func NewDIAdapter(logger *zap.Logger) *DIAdapter {
	return &DIAdapter{
		logger: logger,
	}
}

// Logf implements di logger adapter.
func (d *DIAdapter) Logf(format string, values ...interface{}) {
	d.logger.Debug(fmt.Sprintf("[DI] "+format, values...))
}
