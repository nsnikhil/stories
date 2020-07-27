package store

import (
	"github.com/nsnikhil/stories/pkg/blog/dto"
	"go.uber.org/zap"
)

type TrieStoriesCache struct {
	logger *zap.Logger
}

func NewTrieStoriesCache(logger *zap.Logger) StoriesCache {
	return &TrieStoriesCache{
		logger: logger,
	}
}

func (tc *TrieStoriesCache) AddStory(story *dto.Story) error {
	return nil
}

func (tc *TrieStoriesCache) GetStoryIDs(query string) ([]string, error) {
	return nil, nil
}
