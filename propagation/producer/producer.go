package producer

import (
	"crypto/rsa"

	"github.com/IBM/sarama"

	U "github.com/off-chain-storage/go-off-chain-storage/utils"
)

func syncWriter(brokerList []string) sarama.SyncProducer {
	config := sarama.NewConfig()

	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10

	producer, err := sarama.NewSyncProducer(brokerList, config)
	U.CheckErr(err)

	return producer
}

// 실제 사용할 API (1. Sending Data, 2. Public Key)
func SyncProducer(filedata []byte, pub *rsa.PublicKey) {
	brokerList := []string{
		"localhost:9092",
		"localhost:9093",
		"localhost:9094",
	}

	x := syncWriter(brokerList)

	x.SendMessage(&sarama.ProducerMessage{
		Topic: "off-chain",
		Value: sarama.ByteEncoder(filedata),
	})
}
