package main

import (
	"context"
	"math/rand"
	"time"

	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/kabelsea-sandbox/slice"
	"github.com/kabelsea-sandbox/slice/bundle/monitoring"
	"github.com/kabelsea-sandbox/slice/pkg/di"
)

var (
	SecondCount = stats.Int64("grpc.io/server/received_bytes_per_rpc", "Total bytes received across all messages per RPC.", stats.UnitBytes)
)

type Metrics struct {
}

func NewMetrics() *Metrics {
	go func() {
		for {
			<-time.After(1 * time.Second)
			s := rand.Int63n(65536)
			stats.RecordWithOptions(
				context.Background(),
				stats.WithTags(
					tag.Upsert(ocgrpc.KeyClientMethod, "test"),
				),
				stats.WithMeasurements(SecondCount.M(s)),
			)
		}
	}()
	return &Metrics{}
}

func (m Metrics) Views() []*view.View {
	return []*view.View{
		{
			Name:        "second_count",
			Description: "Second count.",
			Measure:     SecondCount,
			Aggregation: view.Count(),
		},
	}
}

func main() {
	slice.Run(
		slice.UseWorkerDispatcher(),
		slice.RegisterBundles(
			&monitoring.Bundle{
				MetricsEnabled: true,
			},
		),
		slice.ConfigureContainer(
			di.Provide(NewMetrics, di.As(new(monitoring.MetricViews))),
		),
	)
}
