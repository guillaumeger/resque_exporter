package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config := getConfig()
	red := newRedisClient(config)
	getMetrics(red, config)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9112", nil)
}
