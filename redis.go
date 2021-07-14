package main

import (
	"context"
	"sort"
	"strings"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/xtgo/set"
)

var ctx = context.Background()

func newRedisClient(c config) *redis.Client {
	rc := redis.NewClient(&redis.Options{
		Addr:     c.redisHost + ":" + c.redisPort,
		Password: c.redisPassword,
		DB:       c.redisDB,
	})
	return rc
}

func keyExist(c *redis.Client, ns, k string) bool {
	e, err := c.Exists(ns + ":" + k).Result()
	if err != nil {
		log.Errorf("Error contacting redis: %s", err)
		return false
	}
	if e == 0 {
		log.Debugf("Key %s does not exist.", k)
		return false
	}
	log.Debugf("Key %s exists.", k)
	return true
}

func getSetMembers(c *redis.Client, ns, k string) []string {
	m, err := c.SMembers(ns + ":" + k).Result()
	if err != nil {
		log.Errorf("Error contacting redis: %s", err)
		return nil
	}
	return m
}

func getListLength(c *redis.Client, ns, k string) float64 {
	l, err := c.LLen(ns + ":" + k).Result()
	if err != nil {
		log.Errorf("Error contacting redis: %s", err)
		return 0.0
	}
	return float64(l)
}

func getKeyFloat(c *redis.Client, ns, k string) float64 {
	v, err := c.Get(ns + ":" + k).Float64()
	if err != nil {
		log.Errorf("Error contacting redis: %s, or key %s does not exist.", err, k)
		return 0.0
	}
	return v
}

func getKey(c *redis.Client, ns, k string) string {
	v := c.Get(ns + ":" + k).String()
	return v
}

func sanitizeJobName(j string) string {
	s := strings.Split(strings.ReplaceAll(j, "::", "/"), ":")
	return strings.ReplaceAll(s[3], "/", "::")
}

func getJobList(c *redis.Client, ns string) sort.StringSlice {
	iter := c.Scan(0, "resque:stats:jobs*", 0).Iterator()
	var jobs sort.StringSlice
	for iter.Next() {
		jobs = append(jobs, sanitizeJobName(iter.Val()))
	}
	if err := iter.Err(); err != nil {
		log.Errorf("Error retrieving the list of jobs: %s", err)
		return nil
	}
	sort.Sort(jobs)
	u := set.Uniq(jobs)
	return jobs[:u]
}
