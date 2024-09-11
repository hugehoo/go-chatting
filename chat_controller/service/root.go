package service

import (
	"chat_controller/repository"
	"chat_controller/types/table"
	"encoding/json"
	. "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
)

type Server struct {
}

type Service struct {
	repository    *repository.Repository
	AvgServerList map[string]bool
}

func NewService(repository *repository.Repository) *Service {
	s := &Service{repository: repository, AvgServerList: make(map[string]bool)}
	s.setServerInfo()

	if err := s.repository.Kafka.RegisterTopic("chat-server"); err != nil {
		panic(err)
	} else {
		go s.loopKafka()
	}

	// 위에 loopKafka() 는 비동기로 백그라운드에서 실행하고 여기 메서드는 걍 리턴한다.
	// loopKafka() 내부에서 for문이 무한으로 돌고 있음.
	return s
}

func (s *Service) ResponseLiveServerList() []string {
	var res []string
	for ip, available := range s.AvgServerList {
		if available == true {
			res = append(res, ip)
		}
	}
	return res
}

func (s *Service) setServerInfo() {
	if serverList, err := s.GetAvailableServerList(); err != nil {
		panic(err)
	} else {
		for _, server := range serverList {
			s.AvgServerList[server.IP] = true
		}
	}

}

func (s *Service) GetAvailableServerList() ([]*table.ServerInfo, error) {
	return s.repository.GetAvailableServerList()
}

func (s *Service) loopKafka() {
	for {
		ev := s.repository.Kafka.Poll(100)
		switch event := ev.(type) {
		case *Message:

			// Issue : Status 가 status(소문자 시작) 로 돼있으면 대소문자 구분때문에 브로커에서 넘어오는 값을 읽지 못한다.
			type ServerInfoEvent struct {
				Ip     string
				Status bool
			}

			var decoder ServerInfoEvent
			if err := json.Unmarshal(event.Value, &decoder); err != nil {
				log.Println("Failed to decode event", event.Value)
			} else {
				log.Println("📌", decoder)
				s.AvgServerList[decoder.Ip] = decoder.Status
			}
		case *Error:
			log.Println("Failed To Polling Event", event.Error())
		}
	}
}
