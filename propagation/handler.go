package propagation

import (
	"crypto/rsa"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type BlockData struct {
	Filedata []byte         `json:"filedata"`
	Pub      *rsa.PublicKey `json:"pub"`
}

func ConnectRedis() *redis.Client {
	redisDb := redis.NewClient(&redis.Options{
		Addr:         "172.25.0.2:6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	return redisDb
}
