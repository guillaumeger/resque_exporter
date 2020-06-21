package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config := getConfig()

	red := newRedisClient(config)

	// workers
	var (
		workers        float64
		workingWorkers float64
	)

	if !keyExist(red, config.redisNamespace, "workers") {
		fmt.Println("key does not exist")
		workers = 0.0
		workingWorkers = 0
	} else {
		workersList := getSetMembers(red, config.redisNamespace, "workers")
		workers = float64(len(workersList))
		fmt.Println("workers: ", workers)

		workingWorkers = 0.0
		for _, w := range workersList {
			if keyExist(red, config.redisNamespace, "workers:"+w) {
				workingWorkers++
			}
		}
		fmt.Println("working workers: ", workingWorkers)
	}

	promauto.NewGauge(prometheus.GaugeOpts{
		Name: "resque_workers",
		Help: "Number of workers",
	}).Set(workers)

	prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "resque_workers_working",
		Help: "Number of workers currently working",
	}).Set(workingWorkers)

	if keyExist(red, config.redisNamespace, "queues") {
		queuesList := getSetMembers(red, config.redisNamespace, "queues")
		for i, q := range queuesList {
			qJobs := getListLength(red, config.redisNamespace, "queues:"+string(i+1))
			promauto.NewGauge(prometheus.GaugeOpts{
				Name: "resque_queue_jobs",
				Help: "Number of jobs in queue",
				ConstLabels: prometheus.Labels{
					"queue": q,
				},
			}).Set(qJobs)
		}
	}

	promauto.NewCounterFunc(prometheus.CounterOpts{
		Name: "resque_jobs_processed_total",
		Help: "Total number of processed jobs",
	}, func() float64 {
		return getKeyFloat(red, config.redisNamespace, "stat:processed")
	})

	promauto.NewCounterFunc(prometheus.CounterOpts{
		Name: "resque_jobs_failed_total",
		Help: "Total number of failed jobs",
	}, func() float64 {
		return getKeyFloat(red, config.redisNamespace, "stat:failed")
	})

	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "resque_failed_queue",
		Help: "Number of jobs in the failed queue",
	}, func() float64 {
		return getListLength(red, config.redisNamespace, "failed")
	})

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9112", nil)
}
