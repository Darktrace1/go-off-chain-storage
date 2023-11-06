package propagation

import (
	"context"
	"crypto/rsa"

	U "github.com/off-chain-storage/go-off-chain-storage/utils"
	P "github.com/off-chain-storage/go-off-chain-stroage/propagation"
)

// 실제 사용할 API (1. Sending Data, 2. Public Key)
func SyncProducer(filedata []byte, pub *rsa.PublicKey, partitionNum int32) {
	ctx := context.Background()

	// Get Redis Instance
	redisDb := P.ConnectRedis()

	var err error
	err = redisDb.Publish(ctx, "mychannel1", filedata).Err()
	U.CheckErr(err)
}
