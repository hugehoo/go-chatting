package kafka

import (
	"chat_consumer/config"
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Kafka struct {
	Cfg      *config.Config
	Consumer *kafka.Consumer
}
type Message struct {
	Content string `json:"content"`
}

func NewKafka(c *config.Config) (*Kafka, error) {
	k := &Kafka{Cfg: c}
	var err error
	if k.Consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": c.Kafka.Url,
		"group.id":          c.Kafka.ClientId,
		"auto.offset.reset": "latest",
	}); err != nil {
		return nil, err
	} else {
		return k, nil
	}
}

func ConsumeKafka(r *kafka.Consumer, collection *mongo.Collection) {
	for {
		ev := r.Poll(500)
		switch event := ev.(type) {
		case *kafka.Message:
			message := Message{
				Content: string(event.Value),
			}
			log.Println("consuming", message)

			_, err := collection.InsertOne(context.Background(), message)
			if err != nil {
				log.Printf("Error inserting message to MongoDB: %v", err)
			} else {
				log.Printf("Message saved to MongoDB: %s", message.Content)
			}
		case kafka.Error:
			log.Println("Failed To Polling Event", event.Error())
		default:
			log.Println("event:", event)
		}
	}
}
