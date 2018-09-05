package behold

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

// Shared stats
var (
	mHits       = stats.Int64("hits", "The number of hits recieved on / endpoint", "1")
	metricViews = append(append(append(
		ViewsForCounters(mHits),
		ViewsForSums()...),
		ViewsForDistributions()...),
		ViewsForLastValues()...)
)

// Register adds metrics to the exporter
func Register(views ...*view.View) error {
	return view.Register(views...)
}

// ViewsForCounters aggregates defined measures into counters which can be
// registered with Register. A count is a monotonically increasing counter like
// the number of hits on a website.
func ViewsForCounters(mm ...stats.Measure) []*view.View {
	var vv []*view.View
	for _, m := range mm {
		vv = append(vv, &view.View{
			Measure:     m,
			Name:        m.Name(),
			Description: m.Description(),
			Aggregation: view.Count(),
		})
	}
	return vv
}

// ViewsForSums aggregates defined measures into sums which can be registered
// with Register. A sum is a monotonically increasing amount like the amount of
// bytes sent over the wire.
func ViewsForSums(mm ...stats.Measure) []*view.View {
	var vv []*view.View
	for _, m := range mm {
		vv = append(vv, &view.View{
			Measure:     m,
			Name:        m.Name(),
			Description: m.Description(),
			Aggregation: view.Sum(),
		})
	}
	return vv
}

// ViewsForDistributions aggregates defined measures into distributions which
// can be registered with Register. A distribution is a value that should be
// aggregated with a histogram like the latency of requests to an external
// service or the value of a transaction.
func ViewsForDistributions(mm ...stats.Measure) []*view.View {
	var vv []*view.View
	for _, m := range mm {
		vv = append(vv, &view.View{
			Measure:     m,
			Name:        m.Name(),
			Description: m.Description(),
			Aggregation: view.Distribution(),
		})
	}
	return vv
}

// ViewsForLastValues aggregates defined measures into last values which can be
// registered with Register. A last value is a value that changes over time but
// the latest reported value should always be reported. This applies to custom
// aggregations done in the application (e.g. number of goroutines).
func ViewsForLastValues(mm ...stats.Measure) []*view.View {
	var vv []*view.View
	for _, m := range mm {
		vv = append(vv, &view.View{
			Measure:     m,
			Name:        m.Name(),
			Description: m.Description(),
			Aggregation: view.LastValue(),
		})
	}
	return vv
}

// Metrics exposes a prometheus compatible endpoint for metrics collection.
func Metrics(logger *log.Logger, addr string) error {
	p, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		return fmt.Errorf("error creating Prometheus exporter: %v", err)
	}
	view.RegisterExporter(p)
	view.SetReportingPeriod(1 * time.Second)

	if err := Register(metricViews...); err != nil {
		return fmt.Errorf("error registering Prometheus metrics: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", p)
	go func() {
		logger.Printf("serving %s/metrics", addr)
		logger.Fatal(http.ListenAndServe(":8081", mux))
	}()
	return nil
}
