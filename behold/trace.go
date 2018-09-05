package behold

import (
	"fmt"

	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

// Traces initializes the export of traces to jaeger. agent is the UDP address
// of a jaeger agent, endpoint is the URL of the jaeger collector, and service
// is the name and version of the service that is being traced (e.g.
// golo-1.11.6).
func Traces(agent, endpoint, service string) error {
	j, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: agent,
		Endpoint:      endpoint,
		ServiceName:   service,
	})
	if err != nil {
		return fmt.Errorf("error creating Jaeger exporter: %v", err)
	}
	trace.RegisterExporter(j)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	return nil
}

var StartSpan = trace.StartSpan
