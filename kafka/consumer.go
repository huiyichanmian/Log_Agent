package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func Consumer() {
	fmt.Printf("consumer_test")



	// consumer
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Printf("consumer_test create consumer error %s\n", err.Error())
		return
	}

	defer consumer.Close()

	partition_consumer, err := consumer.ConsumePartition("mysql_log", 0, sarama.OffsetOldest)
	if err != nil {
		fmt.Printf("try create partition_consumer error %s\n", err.Error())
		return
	}
	defer partition_consumer.Close()
	fmt.Println(partition_consumer.Messages())
	for {
		select {
		case msg := <-partition_consumer.Messages():
			fmt.Printf("msg offset: %d, partition: %d, timestamp: %s, value: %s\n",
				msg.Offset, msg.Partition, msg.Timestamp.String(), string(msg.Value))
		case err := <-partition_consumer.Errors():
			fmt.Printf("err :%s\n", err.Error())
		}
	}

}
