package service

import (
	"chat_controller/repository"
	"chat_controller/types/table"
)

type Server struct {
}

type Service struct {
	repository    *repository.Repository
	AvgServerList map[string]bool
}

// (s *Service) 얘는 왜 receiver 를 안갖지?
func NewService(repository *repository.Repository) *Service {
	s := &Service{repository: repository, AvgServerList: make(map[string]bool)}
	s.setServerInfo()
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
