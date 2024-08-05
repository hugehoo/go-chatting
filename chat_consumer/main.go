package main

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Message struct {
	Content string `json:"content"`
}

var collection *mongo.Collection

func main() {
	// MongoDB 연결
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection = client.Database("test").Collection("chat")

	// Kafka 소비자 설정
	r, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "127.0.0.1:9092,127.0.0.1:9093,127.0.0.1:9094",
		"group.id":          "chat-consumer-1",
		"auto.offset.reset": "latest",
	})
	defer r.Close()

	// Subscribe to the Kafka topic
	topics := []string{"chat-message"}
	err = r.SubscribeTopics(topics, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to Kafka topics: %v", err)
	}

	// Kafka 메시지 소비 및 MongoDB 저장을 위한 고루틴 시작
	go consumeKafka(r)

	// Fiber 앱 설정
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Kafka consumer is running")
	})

	log.Fatal(app.Listen(":4000"))
}

func consumeKafka(r *kafka.Consumer) {
	for {
		ev := r.Poll(100)
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
