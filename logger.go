package slice

import (
	"fmt"
	"log"
	"strings"
)

// Logger logs slice system messages.
type Logger interface {
	// Debugf logs debug information of component.
	Debugf(component string, format string, values ...interface{})
	// Infof logs information of component.
	Infof(component string, format string, values ...interface{})
	// Errorf logs errors of components.
	Errorf(component string, format string, values ...interface{})
	// Fatalf
	Fatalf(component string, format string, values ...interface{})
}

var (
// default log
// stdLog = &stdLogger{}

// nop log
// nopLog = &errorLogger{}
)

type stdLogger struct {
}

func (l stdLogger) Debugf(component string, format string, values ...interface{}) {
	log.Printf("[%s] %s", strings.ToUpper(component), fmt.Sprintf(format, values...))
}

func (l stdLogger) Infof(component string, format string, values ...interface{}) {
	log.Printf("[%s] %s", strings.ToUpper(component), fmt.Sprintf(format, values...))
}

func (l stdLogger) Errorf(component string, format string, values ...interface{}) {
	log.Printf("[%s] Error: %s", strings.ToUpper(component), fmt.Sprintf(format, values...))
}

func (l stdLogger) Fatalf(component string, format string, values ...interface{}) {
	log.Printf("[%s] Error: %s", strings.ToUpper(component), fmt.Sprintf(format, values...))
}

type errorLogger struct{}

func (l errorLogger) Debugf(string, string, ...interface{}) {}
func (l errorLogger) Infof(string, string, ...interface{})  {}
func (l errorLogger) Errorf(component string, format string, values ...interface{}) {
	log.Printf("[%s] Error: %s", strings.ToUpper(component), fmt.Sprintf(format, values...))
}
func (l errorLogger) Fatalf(component string, format string, values ...interface{}) {
	log.Fatalf("[%s] Error: %s", strings.ToUpper(component), fmt.Sprintf(format, values...))
}
