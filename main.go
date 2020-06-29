package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func init() {
	logLevel := getConfigValue("RESQUE_EXPORTER_LOG_LEVEL", "error")
	levels := make(map[string]log.Level)
	levels["debug"] = log.DebugLevel
	levels["info"] = log.InfoLevel
	levels["warn"] = log.WarnLevel
	levels["error"] = log.ErrorLevel
	levels["fatal"] = log.FatalLevel
	log.SetLevel(levels[logLevel])
}

func main() {
	log.Infof("Starting...")
	log.Debugf("Getting configuration")
	config := getConfig()
	log.Debugf("Configuration: %+v", config)
	red := newRedisClient(config)
	log.Debugf("Created redis client: %+v", red)
	go getWorkersMetrics(red, config)
	go getQueuedJobsMetrics(red, config)
	go getProcessedJobsMetrics(red, config)
	go getFailedJobsMetrics(red, config)
	go getFailedQueueMetrics(red, config)
	log.Debugf("Serving /metrics on port 9112")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9447", nil)
}
