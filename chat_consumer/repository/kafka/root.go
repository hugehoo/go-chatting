package kafka

import (
	"chat_consumer/config"
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hamba/avro/v2"
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
type SimpleRecord struct {
	Name    string `avro:"name"`
	Message string `avro:"message"`
	Room    string `avro:"room"`
}

const recordScheme = `{
        "type": "record",
        "name": "simple",
        "namespace": "org.hamba.avro",
        "fields" : [
            {"name": "name", "type": "string"},
            {"name": "message", "type": "string"},
            {"name": "room", "type": "string"}
        ]
    }`

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
	scheme, err := avro.Parse(recordScheme)
	for {
		ev := r.Poll(500)
		switch event := ev.(type) {
		case *kafka.Message:

			// original Json value
			//message := Message{
			//	Content: string(event.Value),
			//}

			out := SimpleRecord{}
			if err = avro.Unmarshal(scheme, event.Value, &out); err != nil {
				log.Fatal(err)
			}
			log.Println("consuming", out)
			_, err := collection.InsertOne(context.Background(), out)
			if err != nil {
				log.Printf("Error inserting message to MongoDB: %v", err)
			} else {
				log.Printf("Message saved to MongoDB: %s", out)
			}
		case kafka.Error:
			log.Println("Failed To Polling Event", event.Error())
		default:
			log.Println("event:", event)
		}
	}
}
