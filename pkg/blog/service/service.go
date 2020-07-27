package service

type Service struct {
	ss StoriesService
}

func NewService(ss StoriesService) *Service {
	return &Service{
		ss: ss,
	}
}

func (s *Service) GetStoriesService() StoriesService {
	return s.ss
}
