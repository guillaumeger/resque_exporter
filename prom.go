package main

import (
	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type Exporter struct {
	red  *redis.Client
	conf config
}

func NewExporter(red *redis.Client, conf config) *Exporter {
	return &Exporter{
		red:  red,
		conf: conf,
	}
}

var workersMetric = prometheus.NewDesc(
	prometheus.BuildFQName("resque", "", "workers"),
	"Number of workers",
	[]string{},
	nil,
)

func getWorkersMetrics(red *redis.Client, conf config) float64 {
	if !keyExist(red, conf.redisNamespace, "workers") {
		log.Warnf("Key %v does not exist in redis, skipping...", conf.redisNamespace+":workers")
		return 0.0
	} else {
		workersList := getSetMembers(red, conf.redisNamespace, "workers")
		workers := float64(len(workersList))
		log.Debugf("No of workers: %v", workers)
		return workers
	}
}

var schedulesMetric = prometheus.NewDesc(
	prometheus.BuildFQName("resque", "", "schedules"),
	"Number of scheduled jobs",
	[]string{},
	nil,
)

func getSchedulerMetrics(red *redis.Client, conf config) float64 {
	if !keyExist(red, conf.redisNamespace, "schedules_changed") {
		log.Warnf("Key %v does not exist in redis, skipping...", conf.redisNamespace+":schedules_changed")
		return 0.0
	} else {
		schedulesList := getSetMembers(red, conf.redisNamespace, "schedules_changed")
		schedules := float64(len(schedulesList))
		log.Debugf("No of schedules: %v", schedules)
		return schedules
	}
}

var workingWorkersMetric = prometheus.NewDesc(
	prometheus.BuildFQName("resque", "", "workers_working"),
	"Number of workers currently working",
	[]string{},
	nil,
)

func getWorkingWorkersMetric(red *redis.Client, conf config) float64 {
	if !keyExist(red, conf.redisNamespace, "workers") {
		log.Warnf("Key %v does not exist in redis, skipping...", conf.redisNamespace+":workers")
		return 0.0
	} else {
		workingWorkers := 0.0
		workersList := getSetMembers(red, conf.redisNamespace, "workers")
		for _, w := range workersList {
			if keyExist(red, conf.redisNamespace, "workers:"+w) {
				workingWorkers++
			}
		}
		log.Debugf("No of working workers: %v", workingWorkers)
		return workingWorkers
	}
}

var queuedJobsMetric = prometheus.NewDesc(
	prometheus.BuildFQName("resque", "", "queue_jobs"),
	"Number of jobs in queue",
	[]string{
		"queue",
	},
	nil,
)

var processedJobsMetric = prometheus.NewDesc(
	prometheus.BuildFQName("resque", "processed_jobs", "total"),
	"Total number of processed jobs",
	[]string{},
	nil,
)

var failedJobsMetric = prometheus.NewDesc(
	prometheus.BuildFQName("resque", "failed_jobs", "total"),
	"Total number of failed jobs",
	[]string{},
	nil,
)

var failedQueueMetric = prometheus.NewDesc(
	prometheus.BuildFQName("resque", "", "failed_queue"),
	"Number of jobs in the failed queue",
	[]string{},
	nil,
)

var jobFailedMetric = prometheus.NewDesc(
	prometheus.BuildFQName("resque", "", "job_failed"),
	"Failed jobs",
	[]string{
		"job",
	},
	nil,
)

var jobEnqueuedMetric = prometheus.NewDesc(
	prometheus.BuildFQName("resque", "", "job_enqueued"),
	"Enqueued jobs",
	[]string{
		"job",
	},
	nil,
)

var jobCompletedMetric = prometheus.NewDesc(
	prometheus.BuildFQName("resque", "", "job_completed"),
	"Completed jobs",
	[]string{
		"job",
	},
	nil,
)

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- workersMetric
	ch <- schedulesMetric
	ch <- workingWorkersMetric
	ch <- queuedJobsMetric
	ch <- processedJobsMetric
	ch <- failedJobsMetric
	ch <- failedQueueMetric
	ch <- jobFailedMetric
	ch <- jobEnqueuedMetric
	ch <- jobCompletedMetric
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		workersMetric,
		prometheus.GaugeValue,
		getWorkersMetrics(e.red, e.conf),
	)
	ch <- prometheus.MustNewConstMetric(
		schedulesMetric,
		prometheus.GaugeValue,
		getSchedulerMetrics(e.red, e.conf),
	)
	ch <- prometheus.MustNewConstMetric(
		workingWorkersMetric,
		prometheus.GaugeValue,
		getWorkingWorkersMetric(e.red, e.conf),
	)
	if keyExist(e.red, e.conf.redisNamespace, "queues") {
		queuesList := getSetMembers(e.red, e.conf.redisNamespace, "queues")
		for _, q := range queuesList {
			ch <- prometheus.MustNewConstMetric(
				queuedJobsMetric,
				prometheus.GaugeValue,
				getListLength(e.red, e.conf.redisNamespace, "queue:"+q),
				q,
			)
		}
	}
	ch <- prometheus.MustNewConstMetric(
		processedJobsMetric,
		prometheus.CounterValue,
		getKeyFloat(e.red, e.conf.redisNamespace, "stat:processed"),
	)
	ch <- prometheus.MustNewConstMetric(
		failedJobsMetric,
		prometheus.CounterValue,
		getKeyFloat(e.red, e.conf.redisNamespace, "stat:failed"),
	)
	ch <- prometheus.MustNewConstMetric(
		failedQueueMetric,
		prometheus.GaugeValue,
		getListLength(e.red, e.conf.redisNamespace, "failed"),
	)
	jobList := getJobList(e.red, e.conf.redisNamespace)
	if e.conf.resqueStatsMetrics {
		for _, j := range jobList {
			var failed float64
			var completed float64
			if keyExist(e.red, e.conf.redisNamespace, "stats:jobs:"+j+":failed") {
				failed = getKeyFloat(e.red, e.conf.redisNamespace, "stats:jobs:"+j+":failed")
			} else {
				failed = 0
			}
			if keyExist(e.red, e.conf.redisNamespace, "stats:jobs:"+j+":completed") {
				completed = getKeyFloat(e.red, e.conf.redisNamespace, "stats:jobs:"+j+":completed")
			} else {
				completed = 0
			}
			ch <- prometheus.MustNewConstMetric(
				jobFailedMetric,
				prometheus.GaugeValue,
				failed,
				j,
			)
			ch <- prometheus.MustNewConstMetric(
				jobCompletedMetric,
				prometheus.GaugeValue,
				completed,
				j,
			)
			ch <- prometheus.MustNewConstMetric(
				jobEnqueuedMetric,
				prometheus.GaugeValue,
				getKeyFloat(e.red, e.conf.redisNamespace, "stats:jobs:"+j+":enqueued"),
				j,
			)
		}
	}
}
