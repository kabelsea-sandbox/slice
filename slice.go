package slice

import (
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"slice/pkg/di"
)

// Debug is a debug flag.
type Debug bool

// Name is the name of the application.
type Name string

// String return service name string representation
func (n Name) String() string {
	return string(n)
}

// Args is a os.Args alias.
type Args []string

const (
	// Default env variable for debug.
	debugVariableName = "DEBUG"
	// Default env variable name for env.
	envVariableName = "ENV"
	// Default start timeout.
	defaultStartTimeout = 10 * time.Second
	// Default stop timeout.
	defaultStopTimeout = 10 * time.Second
)

// Run runs the application with provided diopts.
func Run(options ...Option) {
	var app = Application{}
	for _, opt := range options {
		opt.apply(&app)
	}
	app.Run(os.Args)
}

// Application is a modular application built around dependency injection.
type Application struct {
	// The environment of application. May be dev, test, prod for development, testing and production respectively.
	Env Env
	// The name of application.
	Name string
	// A functional packages of application.
	Bundles []Bundle
	// The start timeout of application.
	StartTimeout time.Duration
	// The stop timeout of application.
	StopTimeout time.Duration

	// private
	di           []di.Option
	dispatchFunc interface{}
}

// Run starts application. The args is a os.Args in common.
func (a *Application) Run(args []string) {
	// name is required option of application
	if a.Name == "" {
		a.Name = "slice-application"
	}
	// provide env
	env := Env(getEnv(envVariableName))
	a.di = append(a.di, di.Provide(func() Env { return env }))
	// provide debug flag
	debug := getEnv(debugVariableName) == "true"
	a.di = append(a.di, di.Provide(func() Debug { return Debug(debug) }))
	// configure start and stop timeouts
	if a.StartTimeout == 0 {
		a.StartTimeout = defaultStartTimeout
	}
	if a.StopTimeout == 0 {
		a.StopTimeout = defaultStopTimeout
	}
	// provide name, args and context
	a.di = append(a.di,
		di.Provide(func() Name { return Name(a.Name) }),
		di.Provide(func() Args { return args }),
		di.Provide(NewContext),
	)
	container, err := initialization(a.di...)
	if err != nil {
		exitError(err)
	}
	if a.dispatchFunc != nil {
		_ = container.Provide(func() invokeDispatcher {
			return invokeDispatcher{
				fn:        a.dispatchFunc,
				container: container,
			}
		}, di.As(new(Dispatcher)))
	}
	sortedBundles, valid := sortBundles(a.Bundles)
	if !valid {
		exitError(errors.New("bundle dependencies are cyclic"))
	}
	a.Bundles = sortedBundles
	bundles := prepareBundles(a.Bundles...)
	// todo: configurable prefix?
	if err := configureBundles(envconfig.Process, "", bundles...); err != nil {
		exitError(err)
	}
	if err := buildBundles(container, bundles...); err != nil {
		exitError(err)
	}
	var logger Logger
	if !container.Has(&logger) && env.IsDevelopment() {
		_ = container.Provide(func() stdLogger { return stdLogger{} }, di.As(new(Logger)))
	}
	if !container.Has(&logger) && !env.IsDevelopment() {
		_ = container.Provide(func() errorLogger { return errorLogger{} }, di.As(new(Logger)))
	}
	if err = container.Resolve(&logger); err != nil {
		exitError(err)
	}
	// log parameters
	logger.Debugf("slice", "Name: %s", a.Name)
	logger.Debugf("slice", "Env: %s", env)
	logger.Debugf("slice", "Run timeout: %s", a.StartTimeout)
	logger.Debugf("slice", "Stop timeout: %s", a.StopTimeout)
	logger.Debugf("slice", "Debug: %t", debug)
	logger.Debugf("slice", "Args: %s", args[1:])
	// boot
	logger.Debugf("slice", "Boot")
	shutdowns, err := boot(a.StartTimeout, container, bundles...)
	if err != nil {
		exitError(err)
	}
	// run
	logger.Debugf("slice", "Run")
	if runErr := drun(container); runErr != nil {
		if err := reverseShutdown(a.StopTimeout, container, shutdowns); err != nil {
			logger.Errorf("slice", err.Error())
		}
		logger.Fatalf("slice", runErr.Error())
	}
	logger.Debugf("slice", "Shutdown")
	if err := reverseShutdown(a.StopTimeout, container, shutdowns); err != nil {
		logger.Errorf("slice", err.Error())
	}
	container.Cleanup()
}
