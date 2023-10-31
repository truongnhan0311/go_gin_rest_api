package services

import (
	"github.com/go-redis/redis/v8"
	"os"
)

func redisConnect() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     os.Getenv("REDIS_SERVER") + ":" + os.Getenv("REDIS_SERVER_PORT"),
		DB:       0,
		Password: "",
	})

	return client
}
