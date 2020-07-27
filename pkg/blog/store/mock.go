package store

import (
	"github.com/nsnikhil/stories/pkg/blog/dto"
	"github.com/stretchr/testify/mock"
)

type MockStoriesStore struct {
	mock.Mock
}

func (mock *MockStoriesStore) AddStory(story *dto.Story) error {
	args := mock.Called(story)
	return args.Error(0)
}

func (mock *MockStoriesStore) GetStories(storyIDs ...string) ([]dto.Story, error) {
	args := mock.Called(storyIDs)
	return args.Get(0).([]dto.Story), args.Error(1)
}

func (mock *MockStoriesStore) GetMostViewsStories(offset, limit int) ([]dto.Story, error) {
	args := mock.Called(offset, limit)
	return args.Get(0).([]dto.Story), args.Error(1)
}

func (mock *MockStoriesStore) GetTopRatedStories(offset, limit int) ([]dto.Story, error) {
	args := mock.Called(offset, limit)
	return args.Get(0).([]dto.Story), args.Error(1)
}

type MockStoriesCache struct {
	mock.Mock
}

func (mock *MockStoriesCache) AddStory(story *dto.Story) error {
	args := mock.Called(story)
	return args.Error(0)
}

func (mock *MockStoriesCache) GetStoryIDs(query string) ([]string, error) {
	args := mock.Called(query)
	return args.Get(0).([]string), args.Error(1)
}
