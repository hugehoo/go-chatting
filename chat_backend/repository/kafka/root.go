package kafka

import (
	"chat_backend/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Kafka struct {
	cfg      *config.Config
	producer *kafka.Producer
}

func NewKafka(c *config.Config) (*Kafka, error) {
	k := &Kafka{cfg: c}
	var err error
	if k.producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": c.Kafka.URL,
		"client.id":         c.Kafka.ClientID,
		"acks":              "all",
	}); err != nil {
		return nil, err
	} else {
		return k, nil
	}
}

func (k *Kafka) PublishEvent(topic string, value []byte, ch chan kafka.Event) (kafka.Event, error) {
	if err := k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value}, ch); err != nil {
		return nil, err
	} else {
		return <-ch, err
	}
}
