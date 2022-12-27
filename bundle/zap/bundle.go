package zap

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/kabelsea-sandbox/slice/pkg/zapsentry"

	ctxzap "github.com/kabelsea-sandbox/slice/pkg/zaplog/ctx"

	"github.com/kabelsea-sandbox/slice/pkg/zaplog"

	"github.com/kabelsea-sandbox/slice/pkg/di"

	"github.com/kabelsea-sandbox/slice"
)

// Bundle integrate zap logging.
type Bundle struct {
	SentryDSN        string  `envconfig:"SENTRY_DSN"`
	SentrySampleRate float32 `envconfig:"SENTRY_SAMPLE_RATE" default:"1"`
	Sampling         struct {
		Initial    int `envconfig:"INITIAL" default:"100"`
		Thereafter int `envconfig:"THEREAFTER" default:"100"`
	} `envconfig:"ZAP_SAMPLING"`
}

// Build implement bundle interface.
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(b.NewLogger)
	builder.Provide(NewDIAdapter, di.As(new(di.Logger)))
	builder.Provide(NewSliceAdapter, di.As(new(slice.Logger)))
}

// Boot implements BootShutdown interface.
func (b *Bundle) Boot(_ context.Context, interactor slice.Container) (err error) {
	return interactor.Invoke(b.SetupSliceContext)
}

// Shutdown implements Bundle interface.
func (b *Bundle) Shutdown(_ context.Context, _ slice.Container) (err error) {
	return nil
}

// NewLogger creates zap logger.
func (b *Bundle) NewLogger(env slice.Env, debug slice.Debug) (*zap.Logger, func(), error) {
	var cores []zapcore.Core
	// level
	level := zapcore.InfoLevel
	if debug {
		level = zapcore.DebugLevel
	}
	// console core
	encoder := zaplog.NewProductionEncoder()
	if env.IsDevelopment() {
		encoder = zaplog.NewDevelopmentEncoder()
		level = zapcore.DebugLevel
	}
	writeSyncer, _, err := zap.Open("stdout")
	if err != nil {
		return nil, nil, errors.Wrap(err, "zap open stdout")
	}
	console := zapcore.NewCore(encoder, writeSyncer, level) // todo: unlock?
	// console log sampling
	console = zapcore.NewSampler(console, time.Second, b.Sampling.Initial, b.Sampling.Thereafter)
	// add console core
	cores = append(cores, console)
	// sentry core
	if b.SentryDSN != "" {
		sentryCfg := zapsentry.Configuration{
			SentryDsn:        b.SentryDSN,
			SentrySampleRate: b.SentrySampleRate,
		}

		sentry, err := sentryCfg.Build()
		if err != nil {
			return nil, nil, errors.Wrap(err, "zap sentry core")
		}
		// add sentry core
		cores = append(cores, sentry)
	}
	// build main core
	core := zapcore.NewTee(
		cores...,
	)
	// create logger
	logger := zap.New(core)
	cleanup := func() {
		_ = logger.Sync()
	}
	return logger, cleanup, nil
}

// SetupSliceContext setups logger into slice context.
func (b *Bundle) SetupSliceContext(ctx *slice.Context, logger *zap.Logger) {
	ctxzap.SetupContext(ctx.Set, logger)
}
