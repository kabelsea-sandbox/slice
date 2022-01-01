package grpc

import (
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
)

// MetricViews integrates with monitoring bundle via monitoring.MetricViews interface.
type MetricViews struct {
}

// NewMetricViews creates new metric views.
func NewMetricViews() *MetricViews {
	return &MetricViews{}
}

// Views implements monitoring.MetricViews.
func (m *MetricViews) Views() []*view.View {
	var views []*view.View
	views = append(views, ocgrpc.DefaultServerViews...)
	views = append(views, ocgrpc.DefaultClientViews...)
	return views
}
