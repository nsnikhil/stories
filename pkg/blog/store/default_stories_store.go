package store

import (
	"github.com/nsnikhil/stories/pkg/blog/dto"
	"go.uber.org/zap"
)

type DefaultStoriesStore struct {
	logger *zap.Logger
}

func NewDefaultStoriesStore(logger *zap.Logger) StoriesStore {
	return &DefaultStoriesStore{
		logger: logger,
	}
}

func (dss *DefaultStoriesStore) AddStory(story *dto.Story) error {
	return nil
}

func (dss *DefaultStoriesStore) GetStories(storyIDs ...string) ([]dto.Story, error) {
	return nil, nil
}

func (dss *DefaultStoriesStore) GetMostViewsStories(offset, limit int) ([]dto.Story, error) {
	return nil, nil
}

func (dss *DefaultStoriesStore) GetTopRatedStories(offset, limit int) ([]dto.Story, error) {
	return nil, nil
}
