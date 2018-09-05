package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/voutasaurus/octest/behold"
	"github.com/voutasaurus/octest/config"
)

var (
	mHits = behold.NewCounter("hello/hits", "The number of hits recieved on / endpoint", "1")
)

func main() {
	logger := log.New(os.Stderr, "hello: ", log.Llongfile|log.Lmicroseconds|log.LstdFlags)
	logger.Println("starting...")

	var (
		agent       = config.Env("JAEGER_AGENT_ADDR").WithDefault("localhost:6831")
		endpoint    = config.Env("JAEGER_COLLECTOR_URL").WithDefault("http://localhost:14268")
		metricsAddr = config.Env("METRICS_ADDR").WithDefault(":8081")
		addr        = config.Env("HELLO_ADDR").WithDefault(":8080")
	)

	if err := behold.Traces(agent, endpoint, "hello"); err != nil {
		logger.Fatal(err)
	}
	logger.Println("tracing on")
	if err := behold.Metrics(logger, metricsAddr); err != nil {
		logger.Fatal(err)
	}
	logger.Println("metrics on")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, sp := behold.StartSpan(r.Context(), "hello")
		defer sp.End()
		mHits.Record(ctx, 1)
		logger.Println("hit")
		fmt.Fprintln(w, "hello")
	})

	logger.Println("serving on", addr)
	logger.Fatal(http.ListenAndServe(addr, mux))
}
