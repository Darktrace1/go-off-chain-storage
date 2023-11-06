package propagation

import (
	"time"

	redis "github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	redisDb := redis.NewClient(&redis.Options{
		Addr:         ":6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	return redisDb
}
