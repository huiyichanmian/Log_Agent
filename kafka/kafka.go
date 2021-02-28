package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

type logData struct{
	topic string
	data string
}


var (
	client sarama.SyncProducer  // 声明一个全局的连接Kafka的生产者client
	logDataChan chan *logData  // 存放日志的通道
)

func Init(addrs []string,maxSize int)(err error){
	config := sarama.NewConfig()  // 生产者配置
	config.Producer.RequiredAcks = sarama.WaitForAll  // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true  // 成功交付的消息，将在success channel返回

	// 连接kafka
	if client, err = sarama.NewSyncProducer(addrs, config); err != nil{
		fmt.Println("producer closed, err:", err)
		return err
	}
	// 初始化logDataChan
	logDataChan = make(chan *logData,maxSize)
	// 开启后台的goroutine从通道中读取数据发给kafka
	go sendTokafka()
	return nil
}

// 给外部暴露一个函数，该函数只吧日志数据发送到一个内部的channel中
func SendToChan(topic, data string){
	msg := &logData{
		topic:topic,
		data :data,
	}
	logDataChan <- msg
}

// 真正往Kafka发送日志的函数
func sendTokafka(){
	for{
		select {
		case ld := <- logDataChan:
			// 构造一个消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)
			// 发送消息到kafka
			pid, offset, err := client.SendMessage(msg)
			if err != nil{
				fmt.Println("send msg failed, err:",err)
				return
			}
			fmt.Printf("pid:%v offset:%v\n",pid, offset)
		default:
			time.Sleep(time.Millisecond)
		}
	}
}