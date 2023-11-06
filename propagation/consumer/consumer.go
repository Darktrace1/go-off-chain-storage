package consumer

import (
	"context"

	redis "github.com/redis/go-redis/v9"

	U "github.com/off-chain-storage/go-off-chain-storage/utils"
	P "github.com/off-chain-storage/go-off-chain-stroage/propagation"
)

func SyncConsumer() <-chan *redis.Message {
	ctx := context.Background()

	// Get Redis Instance
	redisDb := P.ConnectRedis()

	pubsub := redisDb.Subscribe(ctx, "mychannel1")

	_, err := pubsub.Receive(ctx)
	U.CheckErr(err)

	ch := pubsub.Channel()

	return ch
}
