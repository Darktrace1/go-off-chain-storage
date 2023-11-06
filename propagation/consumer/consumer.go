package propagation

import (
	"fmt"

	"github.com/IBM/sarama"

	U "github.com/off-chain-storage/go-off-chain-storage/utils"
)

func syncListener() sarama.Client {
	config := sarama.NewConfig()

	config.ChannelBufferSize = 1000000

	client, err := sarama.NewClient([]string{"localhost:9092"}, config)
	U.CheckErr(err)

	return client
}

func SyncConsumer(partition_num int32) sarama.PartitionConsumer {
	client := syncListener()

	lastoffset, err := client.GetOffset("off-chain", partition_num, sarama.OffsetNewest)
	U.CheckErr(err)

	consumer, err := sarama.NewConsumerFromClient(client)
	U.CheckErr(err)

	defer func() {
		if err := consumer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("off-chain", partition_num, lastoffset)
	U.CheckErr(err)

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	return partitionConsumer
	// Trap SIGINT to trigger a shutdown.
	// consumed := 0

	// for {
	// 	select {
	// 	case msg := <-partitionConsumer.Messages():
	// 		fmt.Printf("Topic %s Consumed message offset %d , Partition %d\n", msg.Topic, msg.Offset, msg.Partition)
	// 		consumed++
	// 		fmt.Printf("Consumed: %d\n", consumed)
	// 		fmt.Println(string(msg.Key))
	// 		fmt.Println(string(msg.Value))
	// 		fmt.Println("")
	// 	}
	// }
}
