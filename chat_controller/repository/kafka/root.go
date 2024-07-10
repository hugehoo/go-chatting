package kafka

import (
	"chat_controller/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Kafka struct {
	cfg      *config.Config
	consumer *kafka.Consumer
}

func NewKafka(cfg *config.Config) (*Kafka, error) {
	k := &Kafka{cfg: cfg}

	var err error
	if k.consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Url,
		"group.id":          cfg.Kafka.GroupId,
		"auto.offset.reset": "latest",
	}); err != nil {
		panic(err)
	} else {
		return k, nil
	}
}

func (k *Kafka) Poll(timeoutMs int) kafka.Event {
	return k.consumer.Poll(timeoutMs)
}

func (k *Kafka) RegisterTopic(topic string) error {
	if err := k.consumer.Subscribe(topic, nil); err != nil {
		panic(nil)
	}
	return nil
}
