package behold

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

// Metrics exposes a prometheus compatible endpoint for metrics collection.
func Metrics(logger *log.Logger, addr string) error {
	p, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		return fmt.Errorf("error creating Prometheus exporter: %v", err)
	}
	view.RegisterExporter(p)
	view.SetReportingPeriod(1 * time.Second)

	mux := http.NewServeMux()
	mux.Handle("/metrics", p)
	go func() {
		logger.Printf("serving %s/metrics", addr)
		logger.Fatal(http.ListenAndServe(addr, mux))
	}()
	return nil
}

// TODO: provide types for float64 measurements (if necessary)

// Counter is a monotonically increasing counter like the number of hits on a
// website.
type Counter struct {
	*stats.Int64Measure
}

// NewCounter creates and registers a Counter. It will panic if given a name
// that is already in use by another metric.
func NewCounter(name, description, unit string) *Counter {
	m := stats.Int64(name, description, unit)
	if err := view.Register(aggregate(view.Count(), m)); err != nil {
		// trying to register a nil view or one with a name that is
		// already in use is a fatal error. Catch this in CI.
		panic(err)
	}
	return &Counter{m}
}

// Record adds a measurement to the Counter. Optionally pass tags in via ctx.
func (c *Counter) Record(ctx context.Context, v int64) {
	stats.Record(ctx, c.Int64Measure.M(v))
}

// Sum is a monotonically increasing amount like the amount of bytes sent over
// the wire.
type Sum struct {
	*stats.Int64Measure
}

// NewSum creates and registers a Sum. It will panic if given a name that is
// already in use by another metric.
func NewSum(name, description, unit string) *Sum {
	m := stats.Int64(name, description, unit)
	if err := view.Register(aggregate(view.Sum(), m)); err != nil {
		// trying to register a nil view or one with a name that is
		// already in use is a fatal error. Catch this in CI.
		panic(err)
	}
	return &Sum{m}
}

// Record adds a measurement to the Sum. Optionally pass tags in via ctx.
func (s *Sum) Record(ctx context.Context, v int64) {
	stats.Record(ctx, s.Int64Measure.M(v))
}

// Distribution is a value that should be aggregated with a histogram like the
// latency of requests to an external service or the value of a transaction.
type Distribution struct {
	*stats.Int64Measure
}

// NewDistribution creates and registers a Distribution. It will panic if given
// a name that is already in use by another metric.
func NewDistribution(name, description, unit string) *Distribution {
	m := stats.Int64(name, description, unit)
	if err := view.Register(aggregate(view.Distribution(), m)); err != nil {
		// trying to register a nil view or one with a name that is
		// already in use is a fatal error. Catch this in CI.
		panic(err)
	}
	return &Distribution{m}
}

// Record adds a measurement to the Distribution. Optionally pass tags in via ctx.
func (d *Distribution) Record(ctx context.Context, v int64) {
	stats.Record(ctx, d.Int64Measure.M(v))
}

// LastValue is a value that changes over time where the latest value recorded
// by the application is always be reported. This applies to custom
// aggregations done in the application (e.g. number of goroutines).
type LastValue struct {
	*stats.Int64Measure
}

// NewLastValue creates and registers a LastValue. It will panic if given a
// name that is already in use by another metric.
func NewLastValue(name, description, unit string) *LastValue {
	m := stats.Int64(name, description, unit)
	if err := view.Register(aggregate(view.LastValue(), m)); err != nil {
		// trying to register a nil view or one with a name that is
		// already in use is a fatal error. Catch this in CI.
		panic(err)
	}
	return &LastValue{m}
}

// Record adds a measurement to the LastValue. Optionally pass tags in via ctx.
func (lv *LastValue) Record(ctx context.Context, v int64) {
	stats.Record(ctx, lv.Int64Measure.M(v))
}

func aggregate(agg *view.Aggregation, m stats.Measure) *view.View {
	return &view.View{
		Measure:     m,
		Name:        m.Name(),
		Description: m.Description(),
		Aggregation: agg,
	}
}
