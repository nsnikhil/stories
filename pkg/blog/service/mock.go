package service

import (
	"github.com/nsnikhil/stories/pkg/blog/dao"
	"github.com/stretchr/testify/mock"
)

type MockStoriesService struct {
	mock.Mock
}

func (mock *MockStoriesService) AddStory(story *dao.Story) error {
	args := mock.Called(story)
	return args.Error(0)
}

func (mock *MockStoriesService) GetStory(storyID string) (*dao.Story, error) {
	args := mock.Called(storyID)
	return args.Get(0).(*dao.Story), args.Error(1)
}

func (mock *MockStoriesService) SearchStories(query string) ([]dao.Story, error) {
	args := mock.Called(query)
	return args.Get(0).([]dao.Story), args.Error(1)
}

func (mock *MockStoriesService) GetMostViewsStories(offset, limit int) ([]dao.Story, error) {
	args := mock.Called(offset, limit)
	return args.Get(0).([]dao.Story), args.Error(1)
}

func (mock *MockStoriesService) GetTopRatedStories(offset, limit int) ([]dao.Story, error) {
	args := mock.Called(offset, limit)
	return args.Get(0).([]dao.Story), args.Error(1)
}
