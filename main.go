package main

import (
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// var ctx = context.Background()
	red := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	info := red.Info()
	fmt.Println(info)

	keys := red.Keys("resque:*")
	fmt.Println(keys)

	workers, err := red.SMembers("resque:workers").Result()
	if err != nil {
		fmt.Printf("An error occured: %s\n", err)
	}
	workersNo := float64(len(workers))
	fmt.Println(workersNo)

	wNo := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "resque_workers",
		Help: "Number of workers",
	})
	wNo.Set(workersNo)

	wn := 0.0
	for _, w := range workers {
		exist, err := red.Exists("resque:worker" + w).Result()
		if err != nil {
			fmt.Printf("An error occured: %s\n", err)
		}
		if exist == 1 {
			wn++
		}
	}
	fmt.Println(wn)
	prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "resque_workers_working",
		Help: "Number of workers currently working",
	}).Set(wn)

	q, err := red.SMembers("resque:queues").Result()
	if err != nil {
		fmt.Printf("An error occured: %s\n", err)
	}

	fmt.Println(q)

	if len(q) != 0 {
		for _, qName := range q {
			qJobs, err := red.LLen("queue" + qName).Result()
			if err != nil {
				fmt.Printf("An error occured: %s\n", err)
			}
			qJobsNo := promauto.NewGauge(prometheus.GaugeOpts{
				Name: "resque_queue_jobs",
				Help: "Number of jobs in queue",
				ConstLabels: prometheus.Labels{
					"queue": qName,
				},
			})
			qJobsNo.Set(float64(qJobs))
		}
	}

	promauto.NewCounterFunc(prometheus.CounterOpts{
		Name: "resque_jobs_proccesed_total",
		Help: "Total number of proccesed jobs",
	}, func() float64 {
		p, err := red.Get("resque:stat:processed").Float64()
		if err != nil {
			fmt.Printf("An error occured: %s\n", err)
		}
		return p
	})

	promauto.NewCounterFunc(prometheus.CounterOpts{
		Name: "resque_jobs_failed_total",
		Help: "Total number of failed jobs",
	}, func() float64 {
		f, err := red.Get("resque:stat:failed").Float64()
		if err != nil {
			fmt.Printf("An error occured: %s\n", err)
		}
		return f
	})

	fqNo := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "resque_failed_queue",
		Help: "Number of jobs in the failed queue",
	})

	fq, err := red.LLen("resque:failed").Result()
	if err != nil {
		fmt.Printf("An error occured: %s\n", err)
	}

	fqNo.Set(float64(fq))

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9112", nil)
}

func getJobs(c *redis.Client, k string) float64 {
	j, err := c.Get(k).Float64()
	if err != nil {
		fmt.Printf("An error occured: %s\n", err)
	}
	return j
}
