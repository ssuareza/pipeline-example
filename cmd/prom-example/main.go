package main

import (
	"net/http"
	"prom-example/pkg/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Initialize metrics.
	metrics.Init()

	// Create a new ServeMux to use our middleware
	mux := http.NewServeMux()

	// Register handlers with the new mux
	mux.HandleFunc("/", metrics.Health)
	mux.HandleFunc("/test1", metrics.Health)
	mux.HandleFunc("/test2", metrics.Health)
	mux.HandleFunc("/test3", metrics.Health)
	mux.Handle("/metrics", promhttp.Handler())

	// Wrap the entire mux with our instrumentation
	instrumentedHandler := metrics.InstrumentHandler(mux)

	// Start the server with the instrumented handler
	http.ListenAndServe(":2112", instrumentedHandler)
}
