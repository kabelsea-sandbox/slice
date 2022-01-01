package monitoring

import (
	"context"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/pkg/errors"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"

	"github.com/kabelsea-sanbox/slice"
	httpbundle "github.com/kabelsea-sanbox/slice/bundle/http"
	"github.com/kabelsea-sanbox/slice/pkg/di"
)

// Interface helpers for di.As(..)
var (
	IMetricViews = new(MetricViews)
	IHealthCheck = new(HealthCheck)
)

// HealthCheck
type HealthCheck interface {
}

// MetricViews contains opencensus metric view.
type MetricViews interface {
	Views() []*view.View
}

// Bundle is a bundle that provides configured monitoring.
type Bundle struct {
	MetricsEnabled bool `envconfig:"metrics_enabled"`
	TraceEnabled   bool `envconfig:"trace_enabled"`
}

// DependOn
func (b *Bundle) DependOn() []slice.Bundle {
	return []slice.Bundle{
		&httpbundle.Bundle{},
	}
}

// Build provides exporters and worker to di container.
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(NewHealthController, di.As(new(httpbundle.Controller)))

	if b.MetricsEnabled {
		builder.Provide(b.NewPrometheusExporter, di.As(new(view.Exporter)))
		builder.Provide(NewMetricController, di.As(new(httpbundle.Controller)))
	}

	// TODO: skipped ...
	// if b.TraceEnabled {
	// }
}

// Boot registers exporters.
func (b *Bundle) Boot(_ context.Context, container slice.Container) (err error) {
	if b.MetricsEnabled {
		_ = container.Invoke(b.RegisterMetricViews)
		_ = container.Invoke(b.RegisterOpenCensusViewExporter)
	}

	if b.TraceEnabled {
		_ = container.Invoke(b.ConfigureTrace)
		_ = container.Invoke(b.RegisterOpenCensusTraceExporter)
	}
	return nil
}

// Shutdown implements Bundle interface.
func (b *Bundle) Shutdown(_ context.Context, _ slice.Container) (err error) {
	return nil
}

// NewPrometheusExporter creates prometheus exporter.
func (b *Bundle) NewPrometheusExporter() (*prometheus.Exporter, error) {
	exporter, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		return nil, errors.Wrap(err, "prometheus exporter")
	}
	return exporter, nil
}

// RegisterOpenCensusViewExporter registers opencensus view exporter.
func (b *Bundle) RegisterOpenCensusViewExporter(exporter view.Exporter) {
	view.RegisterExporter(exporter)
}

// RegisterOpenCensusTraceExporter registers opencensus trace exporter.
func (b *Bundle) RegisterOpenCensusTraceExporter(exporter trace.Exporter) {
	trace.RegisterExporter(exporter)
}

// ConfigureTrace configures opencensus trace.
func (b *Bundle) ConfigureTrace() {
	trace.ApplyConfig(
		trace.Config{
			DefaultSampler: trace.ProbabilitySampler(0.1),
		},
	)
}

// RegisterMetricViews registers metric views.
func (b *Bundle) RegisterMetricViews(metrics []MetricViews) error {
	var views []*view.View
	for _, metric := range metrics {
		views = append(views, metric.Views()...)
	}
	return view.Register(views...)
}
