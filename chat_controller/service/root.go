package service

type Server struct {
}

type Service struct {
	AvgServerList map[string]bool
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
