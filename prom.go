package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func getWorkersMetrics(red *redis.Client, conf config) {

	workersMetric := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "resque_workers",
		Help: "Number of workers",
	})

	workingWorkersMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "resque_workers_working",
		Help: "Number of workers currently working",
	})

	for {
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
			workers = float64(len(workersList))
			fmt.Println("workers: ", workers)

			workingWorkers = 0.0
			for _, w := range workersList {
				if keyExist(red, conf.redisNamespace, "workers:"+w) {
					workingWorkers++
				}
			}
			fmt.Println("working workers: ", workingWorkers)
		}
		workersMetric.Set(workers)
		workingWorkersMetric.Set(workingWorkers)
		time.Sleep(time.Second * 5)
	}
}

func getQueuedJobsMetrics(red *redis.Client, conf config) {
	queuedJobsMetric := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "resque_queue_jobs",
		Help: "Number of jobs in queue",
	},
		[]string{
			"queue",
		})
	for {
		if keyExist(red, conf.redisNamespace, "queues") {
			queuesList := getSetMembers(red, conf.redisNamespace, "queues")
			for i, q := range queuesList {
				queuedJobsMetric.WithLabelValues(q).Set(getListLength(red, conf.redisNamespace, "queues:"+string(i+1)))
			}
		}
		time.Sleep(time.Second * 5)
	}
}

func getProcessedJobsMetrics(red *redis.Client, conf config) {
	processedJobsMetric := promauto.NewCounter(prometheus.CounterOpts{
		Name: "resque_jobs_processed_total",
		Help: "Total number of processed jobs",
	})

	var origValue float64
	for {
		fmt.Println("original value:", origValue)
		newValue := getKeyFloat(red, conf.redisNamespace, "stat:processed")
		fmt.Println("New value:", newValue)
		diff := newValue - origValue
		fmt.Println("Diff:", diff)
		fmt.Println("-----")

		processedJobsMetric.Add(diff)
		origValue = newValue
		time.Sleep(time.Second * 5)
	}
}

func getFailedJobsMetrics(red *redis.Client, conf config) {
	failedJobsMetric := promauto.NewCounter(prometheus.CounterOpts{
		Name: "resque_jobs_failed_total",
		Help: "Total number of failed jobs",
	})

	var origValue float64
	for {
		fmt.Println("original value:", origValue)
		newValue := getKeyFloat(red, conf.redisNamespace, "stat:failed")
		fmt.Println("New value:", newValue)
		diff := newValue - origValue
		fmt.Println("Diff:", diff)
		fmt.Println("-----")

		failedJobsMetric.Add(diff)
		origValue = newValue
		time.Sleep(time.Second * 5)
	}
}

func getFailedQueueMetrics(red *redis.Client, conf config) {
	failedQueueMetric := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "resque_failed_queue",
		Help: "Number of jobs in the failed queue",
	})

	for {
		failedQueueMetric.Set(getListLength(red, conf.redisNamespace, "failed"))
		time.Sleep(time.Second * 5)
	}
}
