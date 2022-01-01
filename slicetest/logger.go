package slicetest

import (
	"fmt"
	"log"
	"strings"
)

// Logger
type Logger struct {
}

// NewLogger
func NewLogger() *Logger {
	return &Logger{}
}

func (l Logger) Debugf(component string, format string, values ...interface{}) {
	log.Printf("[%s] %s", strings.ToUpper(component), fmt.Sprintf(format, values...))
}

func (l Logger) Infof(component string, format string, values ...interface{}) {
	log.Printf("[%s] %s", strings.ToUpper(component), fmt.Sprintf(format, values...))
}

func (l Logger) Errorf(component string, format string, values ...interface{}) {
	log.Printf("[%s] Error: %s", strings.ToUpper(component), fmt.Sprintf(format, values...))
}

func (l Logger) Fatalf(component string, format string, values ...interface{}) {
	log.Printf("[%s] Error: %s", strings.ToUpper(component), fmt.Sprintf(format, values...))
}
