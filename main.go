package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config := getConfig()
	red := newRedisClient(config)
	go getWorkersMetrics(red, config)
	go getQueuedJobsMetrics(red, config)
	go getProcessedJobsMetrics(red, config)
	go getFailedJobsMetrics(red, config)
	go getFailedQueueMetrics(red, config)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9112", nil)
}
