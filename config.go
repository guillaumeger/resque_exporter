package main

import (
	"fmt"
	"os"
	"strconv"
)

type config struct {
	redisHost      string
	redisPort      string
	redisDB        int
	redisPassword  string
	redisNamespace string
}

func getConfigValue(env, def string) string {
	v, ok := os.LookupEnv(env)
	if !ok {
		return def
	}
	return v
}

func getDBConfig(env string, def int) int {
	v, ok := os.LookupEnv(env)
	if !ok {
		return def
	}
	db, err := strconv.Atoi(v)
	if err != nil {
		fmt.Printf("An error occured: %v\n", err)
	}
	return db
}

func getConfig() config {
	prefix := "RESQUE_EXPORTER"
	configSet := config{
		redisHost:      getConfigValue(prefix+"REDIS_HOST_", "localhost"),
		redisPort:      getConfigValue(prefix+"REDIS_PORT_", "6379"),
		redisPassword:  getConfigValue(prefix+"REDIS_PASSWORD_", ""),
		redisNamespace: getConfigValue(prefix+"REDIS_NAMESPACE_", "resque"),
	}

	configSet.redisDB = getDBConfig(prefix+"REDIS_DB", 0)
	return configSet
}
