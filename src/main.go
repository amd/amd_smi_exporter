package main

import (
	"log"
	"net/http"
	"src/collect" // This has the implementation of the Scan() function
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Called on each collector.Collect.
	handle := func() (collect.AMDParams) {
		return collect.Scan()
	}

	// Make Prometheus client aware of our collector.
	c := NewCollector(handle)
	prometheus.MustRegister(c)

	// Set up HTTP handler for metrics.
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	// Start listening for HTTP connections.
	const addr = ":2021"
	log.Printf("starting collector exporter on %q", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("cannot start collector exporter: %s", err)
	}
}
