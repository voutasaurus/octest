package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

var (
	mHits = stats.Int64("hello/hits", "The number of hits recieved on / endpoint", "1")
	vHits = &view.View{
		Measure:     mHits,
		Name:        mHits.Name(),
		Description: mHits.Description(),
		Aggregation: view.Count(),
	}
)

func main() {
	logger := log.New(os.Stderr, "hello: ", log.Llongfile|log.Lmicroseconds|log.LstdFlags)
	logger.Println("starting...")

	if err := startTracing(); err != nil {
		logger.Fatal(err)
	}
	logger.Println("tracing on")
	if err := startMetrics(); err != nil {
		logger.Fatal(err)
	}
	logger.Println("metrics on")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, sp := trace.StartSpan(r.Context(), "hello")
		defer sp.End()
		stats.Record(ctx, mHits.M(1))
		logger.Println("hit")
		fmt.Fprintln(w, "hello")
	})

	logger.Println("serving on :8080")
	logger.Fatal(http.ListenAndServe(":8080", mux))
}

func startTracing() error {
	j, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: "localhost:9411",
		ServiceName:   "hello",
	})
	if err != nil {
		return fmt.Errorf("error creating Jaeger exporter: %v", err)
	}
	trace.RegisterExporter(j)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	return nil
}

func startMetrics() error {
	p, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		return fmt.Errorf("error creating Prometheus exporter: %v", err)
	}
	view.RegisterExporter(p)
	view.SetReportingPeriod(1 * time.Second)
	view.Register(vHits)
	mux := http.NewServeMux()
	mux.Handle("/metrics", p)
	go func() {
		log.Println("serving :8081/metrics")
		log.Fatal(http.ListenAndServe(":8081", mux))
	}()
	return nil
}
