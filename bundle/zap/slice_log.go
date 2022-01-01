package zap

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// SliceAdapter implements slice.Logger.
type SliceAdapter struct {
	logger *zap.Logger
}

// NewSliceAdapter creates zap implementation for slice.Logger.
func NewSliceAdapter(logger *zap.Logger) *SliceAdapter {
	return &SliceAdapter{logger: logger}
}

// Debugf implements slice.Logger interface.
func (s *SliceAdapter) Debugf(component string, format string, values ...interface{}) {
	s.logger.Debug(fmt.Sprintf("[%s] %s", strings.ToUpper(component), fmt.Sprintf(format, values...)))
}

// Infof implements slice.Logger interface.
func (s *SliceAdapter) Infof(component string, format string, values ...interface{}) {
	s.logger.Info(fmt.Sprintf("[%s] %s", strings.ToUpper(component), fmt.Sprintf(format, values...)))
}

// Errorf implements slice.Logger interface.
func (s *SliceAdapter) Errorf(component string, format string, values ...interface{}) {
	s.logger.Error(fmt.Sprintf("[%s] %s", strings.ToUpper(component), fmt.Sprintf(format, values...)))
}

// Fatalf implements slice.Logger interface.
func (s *SliceAdapter) Fatalf(component string, format string, values ...interface{}) {
	s.logger.Fatal(fmt.Sprintf("[%s] %s", strings.ToUpper(component), fmt.Sprintf(format, values...)))
}
