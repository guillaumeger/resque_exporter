package main

import (
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

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
		return false
	}
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
		log.Errorf("Error contacting redis: %s", err)
		return 0.0
	}
	return v
}
