package main

import (
	"chat_consumer/config"
	messageBroker "chat_consumer/repository/kafka"
	"context"
	"flag"
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
var pathFlag = flag.String("config", "./config.toml", "config set up")

func main() {

	flag.Parse()
	c := config.NewConfig(*pathFlag)

	// MongoDB 연결
	client, err := mongo.NewClient(options.Client().ApplyURI(c.Mongo.Url))
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

	collection = client.Database(c.Mongo.Database).Collection("chat")

	// Kafka 소비자 설정
	r, err := messageBroker.NewKafka(c)
	defer r.Consumer.Close()

	// Subscribe to the Kafka topic
	topics := []string{"chat-message"}
	err = r.Consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to Kafka topics: %v", err)
	}

	// Kafka 메시지 소비 및 MongoDB 저장을 위한 고루틴 시작
	//go consumeKafka(r)
	go messageBroker.ConsumeKafka(r.Consumer, collection)

	// Fiber 앱 설정
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Kafka consumer is running")
	})

	log.Fatal(app.Listen(":4000"))
}
