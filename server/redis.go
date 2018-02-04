package main

import (
	"github.com/go-redis/redis"
)

// StatusKey to save status
const StatusKey string = "Status"

func newRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func redisGetStatus() (string, error) {
	client := newRedisClient()
	val, err := client.Get(StatusKey).Result()
	if err != nil {
		if err == redis.Nil {
			logger.Printf("no record with key %s found, return %s", StatusKey, StatusNormal)
			return StatusNormal, nil
		}
		logger.Printf("retrieve status from redis, %s", err)
		return "", err
	}

	return val, nil
}

func redisSetStatus(status string) error {
	client := newRedisClient()
	err := client.Set(StatusKey, status, 0).Err()
	if err != nil {
		logger.Printf("set status in redis, %s", err)
	}
	return err
}
