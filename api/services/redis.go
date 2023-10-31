package services

import (
	"github.com/chenyahui/gin-cache/persist"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
)

func RedisConnect() *redis.Client {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	client := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     os.Getenv("REDIS_SERVER") + ":" + os.Getenv("REDIS_SERVER_PORT"),
		DB:       db,
		Password: "",
	})

	return client
}

func RedisStore() *persist.RedisStore {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	redisStore := persist.NewRedisStore(redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     os.Getenv("REDIS_SERVER") + ":" + os.Getenv("REDIS_SERVER_PORT"),
		DB:       db,
		Password: "",
	}))
	return redisStore
}
