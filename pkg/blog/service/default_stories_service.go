package service

import (
	"github.com/nsnikhil/stories/pkg/blog/domain"
	"github.com/nsnikhil/stories/pkg/blog/store"
	"go.uber.org/zap"
)

type DefaultStoriesService struct {
	store  store.Store
	logger *zap.Logger
}

func NewDefaultStoriesService(logger *zap.Logger, store store.Store) StoriesService {
	return &DefaultStoriesService{
		store:  store,
		logger: logger,
	}
}

func (dss *DefaultStoriesService) AddStory(story *domain.Story) error {
	return nil
}

func (dss *DefaultStoriesService) GetStory(storyID string) (*domain.Story, error) {
	return nil, nil
}

func (dss *DefaultStoriesService) SearchStories(query string) ([]domain.Story, error) {
	return nil, nil
}

func (dss *DefaultStoriesService) GetMostViewsStories(offset, limit int) ([]domain.Story, error) {
	return nil, nil
}

func (dss *DefaultStoriesService) GetTopRatedStories(offset, limit int) ([]domain.Story, error) {
	return nil, nil
}
