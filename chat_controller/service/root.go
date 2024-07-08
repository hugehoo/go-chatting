package service

import "chat_controller/repository"

type Server struct {
}

type Service struct {
	repository    *repository.Repository
	AvgServerList map[string]bool
}

// (s *Service) 얘는 왜 receiver 를 안갖지?
func NewService(repository *repository.Repository) *Service {
	s := &Service{repository: repository, AvgServerList: make(map[string]bool)}
	return s
}

func (s *Service) GetAvgServerList() []string {
	var res []string
	for ip, available := range s.AvgServerList {
		if available == true {
			res = append(res, ip)
		}
	}
	return res
}
