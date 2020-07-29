package service

import (
	"github.com/nsnikhil/stories/pkg/blog/domain"
	"github.com/stretchr/testify/mock"
)

type MockStoriesService struct {
	mock.Mock
}

func (mock *MockStoriesService) AddStory(story *domain.Story) error {
	args := mock.Called(story)
	return args.Error(0)
}

func (mock *MockStoriesService) GetStory(storyID string) (*domain.Story, error) {
	args := mock.Called(storyID)
	return args.Get(0).(*domain.Story), args.Error(1)
}

func (mock *MockStoriesService) UpdateStory(story *domain.Story) (int64, error) {
	args := mock.Called(story)
	return args.Get(0).(int64), args.Error(1)
}

func (mock *MockStoriesService) SearchStories(query string) ([]domain.Story, error) {
	args := mock.Called(query)
	return args.Get(0).([]domain.Story), args.Error(1)
}

func (mock *MockStoriesService) GetMostViewsStories(offset, limit int) ([]domain.Story, error) {
	args := mock.Called(offset, limit)
	return args.Get(0).([]domain.Story), args.Error(1)
}

func (mock *MockStoriesService) GetTopRatedStories(offset, limit int) ([]domain.Story, error) {
	args := mock.Called(offset, limit)
	return args.Get(0).([]domain.Story), args.Error(1)
}
