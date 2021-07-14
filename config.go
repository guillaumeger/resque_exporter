package main

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type config struct {
	redisHost          string
	redisPort          string
	redisDB            int
	redisPassword      string
	redisNamespace     string
	resqueStatsMetrics bool
}

func getConfigValue(env, def string) string {
	v, ok := os.LookupEnv(env)
	if !ok {
		return def
	}
	return v
}

func getConfigValueBool(env string, def bool) bool {
	v, ok := os.LookupEnv(env)
	if !ok {
		return def
	}
	val, err := strconv.ParseBool(v)
	if err != nil {
		log.Errorf("Error converting the env variable %s to boolean.", v)
	}
	return val
}

func getDBConfig(env string, def int) int {
	v, ok := os.LookupEnv(env)
	if !ok {
		return def
	}
	db, err := strconv.Atoi(v)
	if err != nil {
		log.Errorf("An error occured: %v\n", err)
	}
	return db
}

func getConfig() config {
	prefix := "RESQUE_EXPORTER_"
	configSet := config{
		redisHost:          getConfigValue(prefix+"REDIS_HOST", "localhost"),
		redisPort:          getConfigValue(prefix+"REDIS_PORT", "6379"),
		redisPassword:      getConfigValue(prefix+"REDIS_PASSWORD", ""),
		redisNamespace:     getConfigValue(prefix+"REDIS_NAMESPACE", "resque"),
		resqueStatsMetrics: getConfigValueBool(prefix+"RESQUE_STATS_METRICS", false),
	}

	configSet.redisDB = getDBConfig(prefix+"REDIS_DB", 0)
	return configSet
}
