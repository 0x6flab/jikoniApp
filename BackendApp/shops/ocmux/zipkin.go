package ocmux

import (
	"io"
	"log"
	"time"

	"contrib.go.opencensus.io/exporter/prometheus"
	oczipkin "contrib.go.opencensus.io/exporter/zipkin"
	"contrib.go.opencensus.io/integrations/ocsql"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

const (
	duration = 10 * time.Second
)

// InitOpenCensusWithZipkin initializes the OpenCensus Zipkin Exporter.
func InitOpenCensusWithZipkin(zipkinURL, serviceName, hostPort string) io.Closer {
	// Always sample our traces for this demo.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// The Zipkin reporter takes collected spans from the app and reports them to the backend
	// http://localhost:9411/api/v2/spans is the default for the Zipkin Span v2
	rep := zipkinhttp.NewReporter(zipkinURL)

	// The localEndpoint stores the name and address of the local service
	localEndpoint, err := zipkin.NewEndpoint(serviceName, hostPort)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}
	exporter := oczipkin.NewExporter(rep, localEndpoint)

	// Register our trace exporter.
	trace.RegisterExporter(exporter)

	// Enable ocsql metrics with OpenCensus
	ocsql.RegisterAllViews()

	// set up our Prometheus details
	prometheusExporter, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "ocsql",
	})
	if err != nil {
		log.Fatalf("error configuring prometheus: %v\n", err)
	}

	// add default ochttp server views
	if err := view.Register(ochttp.DefaultServerViews...); err != nil {
		log.Fatalf("error configuring prometheus: %v\n", err)
	}

	// Report stats at every second.
	view.SetReportingPeriod(duration)

	// use Prometheus for metrics
	view.RegisterExporter(prometheusExporter)

	return rep
}
