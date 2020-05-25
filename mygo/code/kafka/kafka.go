package main

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	kafkaHost = "49.235.235.221:9092"
	groupID = "group-1"
	topic = "helloKafka"
)

func Consumer() {
	log.Println("start product ...")

	cConfig := cluster.NewConfig()
	cConfig.Consumer.Return.Errors = true
	cConfig.Group.Return.Notifications = true
	cConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	cConfig.Consumer.Offsets.CommitInterval = time.Second

	consumer, err := cluster.NewConsumer(strings.Split(kafkaHost,","), groupID, strings.Split(topic,","), cConfig)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()
	// consume messages, watch errors and notifications
	for {
		select {
		case msg, more := <-consumer.Messages():
			if more {
				consumer.MarkOffset(msg, "") // mark message as processed
				log.Printf("consumer msg: %s", msg.Value)
			}
		case err, more := <-consumer.Errors():
			if more {
				log.Printf("Error: %s\n", err.Error())
			}
		case ntf, more := <-consumer.Notifications():
			if more {
				log.Printf("Rebalanced: %+v\n", ntf)
			}
		}
	}
}


func Producer() {
	pConfig := sarama.NewConfig()
	pConfig.Producer.Return.Successes = true
	pConfig.Producer.Timeout = 2 * time.Second
	pConfig.Producer.Compression = sarama.CompressionGZIP
	pConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	p, err := sarama.NewAsyncProducer(strings.Split(kafkaHost,","), pConfig)
	if err != nil {
		panic(err)
	}
	defer p.AsyncClose()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder("go_test"),
	}
	var id int
	for {
		id += 1
		t := time.Now().Format("2006-01-02 15:04:05")
		data := map[string]string{
			"time":t,
			"id": strconv.Itoa(id),
		}
		dJson, _ :=json.Marshal(data)
		msg.Value = sarama.ByteEncoder(dJson)
		p.Input() <- msg

		select {
		case suc := <-p.Successes():
			log.Printf("offset: %d,  timestamp: %s", suc.Offset, suc.Timestamp.String())
		case fail := <-p.Errors():
			log.Printf("err: %s\n", fail.Err.Error())
		}

		time.Sleep(1 * time.Second)
	}

}
