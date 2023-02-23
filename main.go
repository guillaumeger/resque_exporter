package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func init() {
	logLevel := getConfigValue("RESQUE_EXPORTER_LOG_LEVEL", "info")
	levels := make(map[string]log.Level)
	levels["debug"] = log.DebugLevel
	levels["info"] = log.InfoLevel
	levels["warn"] = log.WarnLevel
	levels["error"] = log.ErrorLevel
	levels["fatal"] = log.FatalLevel
	log.SetLevel(levels[logLevel])
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

func main() {
	fmt.Println("Starting")
	log.Infof("Starting...")
	log.Debugf("Getting configuration")
	config := getConfig()
	log.Debugf("Configuration: %+v", config)
	red := newRedisClient(config)
	log.Debugf("Created redis client: %+v", red)
	exporter := NewExporter(red, config)
	prometheus.MustRegister(exporter)
	log.Debugf("Serving /metrics on port 9447")
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":9447", nil)
	if err != nil {
		log.Errorf("Cannot serve on port 9447: %s", err)
	}
}
