package monitoring

import (
	"net/http"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/go-chi/chi/v5"
)

const (
	defaultMetricsPath     = "/metrics"
	defaultHealthCheckPath = "/health-check"
)

// HealthCheckController
type HealthCheckController struct {
}

// NewHealthController constructs controller.
func NewHealthController() *HealthCheckController {
	return &HealthCheckController{}
}

// RegisterRoutes
func (c HealthCheckController) RegisterRoutes(mux chi.Router) {
	mux.HandleFunc(defaultHealthCheckPath, func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK) // todo: add real health checks
	})
}

// MetricController
type MetricController struct {
	exporter *prometheus.Exporter
}

// NewMetricController constructs controller.
func NewMetricController(exporter *prometheus.Exporter) *MetricController {
	return &MetricController{exporter: exporter}
}

// RegisterRoutes
func (m MetricController) RegisterRoutes(mux chi.Router) {
	mux.Handle(defaultMetricsPath, m.exporter)
}
