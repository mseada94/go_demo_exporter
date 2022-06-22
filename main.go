package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordCounterMetric() {
	// Change Metric over time
	go func() {
		for {
			counterMetric.Inc()
			time.Sleep(30 * time.Second)
		}
	}()
}

func recordGaugeMetric() {
	// Set Metric once
	gaugeMetric.Set(95)
}

func createGaugeCollector(instance string) prometheus.Collector {
	// Define a callback function to read the Metric value with each request
	return prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "go_demo_gauge_collector",
			Help: "A metric with a random value and constant labels labeled",
			ConstLabels: prometheus.Labels{
				"version":  "1.0",
				"branch":   "dev",
				"instance": instance,
			},
		},
		func() float64 { return rand.Float64() * 100 },
	)
}

var (
	counterMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "go_demo_counter",
		Help: "Counter Metric Type Demo",
	})

	gaugeMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "go_demo_gauge",
		Help: "Gauge Metric Type Demo",
	})
)

func main() {

	recordCounterMetric()
	recordGaugeMetric()

	prometheus.Register(createGaugeCollector("instance-1"))
	prometheus.Register(createGaugeCollector("instance-2"))

	println("Exporter listening on port 2112")
	println("open http://localhost:2112/metrics")

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

}
