package main

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func getMetrics(red *redis.Client, conf config) {
	var (
		workers        float64
		workingWorkers float64
	)

	if !keyExist(red, conf.redisNamespace, "workers") {
		fmt.Println("key does not exist")
		workers = 0.0
		workingWorkers = 0
	} else {
		workersList := getSetMembers(red, conf.redisNamespace, "workers")
		promauto.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "resque_workers",
			Help: "Number of workers",
		}, func() float64 { return float64(len(workersList)) })
		fmt.Println("workers: ", workers)

		workingWorkers = 0.0
		for _, w := range workersList {
			if keyExist(red, conf.redisNamespace, "workers:"+w) {
				workingWorkers++
			}
		}
		fmt.Println("working workers: ", workingWorkers)
	}

	prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "resque_workers_working",
		Help: "Number of workers currently working",
	}, func() float64 { return workingWorkers })

	if keyExist(red, conf.redisNamespace, "queues") {
		queuesList := getSetMembers(red, conf.redisNamespace, "queues")
		for i, q := range queuesList {
			promauto.NewGaugeFunc(prometheus.GaugeOpts{
				Name: "resque_queue_jobs",
				Help: "Number of jobs in queue",
				ConstLabels: prometheus.Labels{
					"queue": q,
				},
			}, func() float64 { return getListLength(red, conf.redisNamespace, "queues:"+string(i+1)) },
			)
		}
	}

	promauto.NewCounterFunc(prometheus.CounterOpts{
		Name: "resque_jobs_processed_total",
		Help: "Total number of processed jobs",
	}, func() float64 {
		return getKeyFloat(red, conf.redisNamespace, "stat:processed")
	})

	promauto.NewCounterFunc(prometheus.CounterOpts{
		Name: "resque_jobs_failed_total",
		Help: "Total number of failed jobs",
	}, func() float64 {
		return getKeyFloat(red, conf.redisNamespace, "stat:failed")
	})

	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "resque_failed_queue",
		Help: "Number of jobs in the failed queue",
	}, func() float64 {
		return getListLength(red, conf.redisNamespace, "failed")
	})
}
