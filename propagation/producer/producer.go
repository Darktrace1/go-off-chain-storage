package producer

import (
	"github.com/IBM/sarama"

	U "github.com/off-chain-storage/go-off-chain-storage/utils"
)

func SyncWriter(brokerList []string) sarama.SyncProducer {
	config := sarama.NewConfig()

	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10

	producer, err := sarama.NewSyncProducer(brokerList, config)
	U.CheckErr(err)

	return producer
}

func SyncProducer(filedata []byte) {
	brokerList := []string{
		"localhost:9092",
		"localhost:9093",
		"localhost:9094",
	}

	x := SyncWriter(brokerList)

	x.SendMessage(&sarama.ProducerMessage{
		Topic: "off-chain",
		Value: sarama.ByteEncoder(filedata),
	})
}
