package main

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

func getWorkersMetrics(red *redis.Client, conf config) {

	workersMetric := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "resque_workers",
		Help: "Number of workers",
	})
	log.Debugf("Created metric: resque_workers")

	workingWorkersMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "resque_workers_working",
		Help: "Number of workers currently working",
	})
	log.Debugf("Created metric: resque_workers_working")
	for {
		log.Debugf("Starting getWorkersMetrics loop")
		var (
			workers        float64
			workingWorkers float64
		)

		if !keyExist(red, conf.redisNamespace, "workers") {
			log.Warnf("Key %v does not exist in redis, skipping...", conf.redisNamespace+":workers")
			workers = 0.0
			workingWorkers = 0.0
		} else {
			workersList := getSetMembers(red, conf.redisNamespace, "workers")
			workers = float64(len(workersList))
			log.Debugf("No of workers: %v", workers)

			workingWorkers = 0.0
			for _, w := range workersList {
				if keyExist(red, conf.redisNamespace, "workers:"+w) {
					workingWorkers++
				}
			}
			log.Debugf("No of working workers: %v", workingWorkers)
		}
		log.Debugf("Setting metric value for resque_workers, value: %v", workers)
		workersMetric.Set(workers)
		log.Debugf("Setting metric value for resque_workers_working, value: %v", workingWorkers)
		workingWorkersMetric.Set(workingWorkers)
		log.Debugf("Sleeping 5 seconds")
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
	log.Debugf("Created metric: resque_queue_jobs")
	for {
		log.Debugf("Starting getQueuedJobsMetrics loop")
		if keyExist(red, conf.redisNamespace, "queues") {
			queuesList := getSetMembers(red, conf.redisNamespace, "queues")
			for i, q := range queuesList {
				log.Debugf("Setting metric value for resque_queue_job")
				queuedJobsMetric.WithLabelValues(q).Set(getListLength(red, conf.redisNamespace, "queues:"+string(i+1)))
			}
		}
		log.Debugf("Sleeping 5 seconds")
		time.Sleep(time.Second * 5)
	}
}

func getProcessedJobsMetrics(red *redis.Client, conf config) {
	processedJobsMetric := promauto.NewCounter(prometheus.CounterOpts{
		Name: "resque_jobs_processed_total",
		Help: "Total number of processed jobs",
	})
	log.Debugf("Created metric: resque_jobs_processed_total")

	var origValue float64
	for {
		log.Debugf("Starting getProcessedJobsMetrics loop")
		log.Debugf("Calculating the difference of processed jobs between loop iterations")
		log.Debugf("Processed jobs original value: %v", origValue)
		newValue := getKeyFloat(red, conf.redisNamespace, "stat:processed")
		log.Debugf("Processed jobs new value: %v", newValue)
		diff := newValue - origValue
		log.Debugf("Processed jobs difference: %v", diff)

		log.Debugf("Setting metric value for resque_jobs_processed_total, adding %v", diff)
		processedJobsMetric.Add(diff)
		origValue = newValue
		log.Debugf("Sleeping 5 seconds")
		time.Sleep(time.Second * 5)
	}
}

func getFailedJobsMetrics(red *redis.Client, conf config) {
	failedJobsMetric := promauto.NewCounter(prometheus.CounterOpts{
		Name: "resque_jobs_failed_total",
		Help: "Total number of failed jobs",
	})
	log.Debugf("Created metric: resque_jobs_failed_total")

	var origValue float64
	for {
		log.Debugf("Starting getFailedJobsMetrics loop")
		log.Debugf("Calculating the difference of failed processed jobs between loop iterations")
		log.Debugf("Failed jobs original value: %v", origValue)
		newValue := getKeyFloat(red, conf.redisNamespace, "stat:failed")
		log.Debugf("Failed jobs failed value: %v", origValue)
		diff := newValue - origValue
		log.Debugf("Failed jobs difference: %v", origValue)

		log.Debugf("Setting metric value for resque_jobs_failed_total, adding %v", diff)
		failedJobsMetric.Add(diff)
		origValue = newValue
		log.Debugf("Sleeping 5 seconds")
		time.Sleep(time.Second * 5)
	}
}

func getFailedQueueMetrics(red *redis.Client, conf config) {
	failedQueueMetric := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "resque_failed_queue",
		Help: "Number of jobs in the failed queue",
	})
	log.Debugf("Created metric: resque_failed_queue")

	for {
		log.Debugf("Starting getFailedQueueMetrics loop")
		log.Debugf("Setting metric value for resque_failed_queue")
		failedQueueMetric.Set(getListLength(red, conf.redisNamespace, "failed"))
		log.Debugf("Sleeping 5 seconds")
		time.Sleep(time.Second * 5)
	}
}
