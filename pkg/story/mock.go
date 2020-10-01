package story

import (
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/stretchr/testify/mock"
)

type MockStoriesService struct {
	mock.Mock
}

func (mock *MockStoriesService) AddStory(story *model.Story) error {
	args := mock.Called(story)
	return args.Error(0)
}

func (mock *MockStoriesService) GetStory(storyID string) (*model.Story, error) {
	args := mock.Called(storyID)
	return args.Get(0).(*model.Story), args.Error(1)
}

func (mock *MockStoriesService) UpdateStory(story *model.Story) (int64, error) {
	args := mock.Called(story)
	return args.Get(0).(int64), args.Error(1)
}

func (mock *MockStoriesService) DeleteStory(storyID string) (int64, error) {
	args := mock.Called(storyID)
	return args.Get(0).(int64), args.Error(1)
}

func (mock *MockStoriesService) SearchStories(query string) ([]model.Story, error) {
	args := mock.Called(query)
	return args.Get(0).([]model.Story), args.Error(1)
}

func (mock *MockStoriesService) GetMostViewsStories(offset, limit int) ([]model.Story, error) {
	args := mock.Called(offset, limit)
	return args.Get(0).([]model.Story), args.Error(1)
}

func (mock *MockStoriesService) GetTopRatedStories(offset, limit int) ([]model.Story, error) {
	args := mock.Called(offset, limit)
	return args.Get(0).([]model.Story), args.Error(1)
}
