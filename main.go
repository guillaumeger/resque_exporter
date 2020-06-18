package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var ctx = context.Background()
	red := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// info := red.Info(ctx)
	// fmt.Println(info)
	keys := red.Keys(ctx, "resque:*")
	fmt.Println(keys)
	workers, err := red.SMembers(ctx, "resque:workers").Result()
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
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9112", nil)
}
