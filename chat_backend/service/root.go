package service

import (
	"chat_backend/repository"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type Service struct {
	repository *repository.Repository
}

func NewService(rep *repository.Repository) *Service {
	return _
}

func (s *Service) ServerSet(ip string, available bool) error {
	if err := s.repository.ServerSet(ip, available); err != nil {
		log.Println("Failed To ServerSet", "ip", ip, "available", available)
		return err
	} else {
		return nil
	}
}

func (s *Service) PublishServerStatusEvent(ip string, status bool) {
	type ServerInfoEvent struct {
		IP     string
		Status bool
	}
	event := &ServerInfoEvent{IP: ip, Status: status}
	ch := make(chan kafka.Event)

	if v, err := json.Marshal(event); err != nil {
		log.Println("Failed to marshal")
	} else if result, err := s.PublishEvent("chat", v, ch); err != nil {
		log.Println("Failed To Send Event To Kafka", "err", err)
	} else {
		log.Println("Success To Send Event", result)
	}
}

func (s *Service) PublishEvent(topic string, value []byte, ch chan kafka.Event) (kafka.Event, error) {
	return s.repository.Kafka.PublishEvent(topic, value, ch)
}
