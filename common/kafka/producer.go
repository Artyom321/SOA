package kafka

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func NewProducer(broker string) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	for range time.Tick(time.Second) {
		producer, err := sarama.NewSyncProducer([]string{broker}, config)
		if err == nil {
			return producer
		}
		log.Println(broker)
		log.Println("Retrying to create broker connection...")
		time.Sleep(time.Second)
	}

	return nil
}

func SendMessage(producer sarama.SyncProducer, topic string, value interface{}) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(valueBytes),
	}

	_, _, err = producer.SendMessage(&msg)
	if err != nil {
		log.Println("Failed to send message to broker")
		return err
	}

	return nil
}
