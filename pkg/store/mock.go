package store

import (
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/stretchr/testify/mock"
)

type MockStoriesStore struct {
	mock.Mock
}

func (mock *MockStoriesStore) AddStory(story *model.Story) error {
	args := mock.Called(story)
	return args.Error(0)
}

func (mock *MockStoriesStore) GetStories(storyIDs ...string) ([]model.Story, error) {
	args := mock.Called(storyIDs)
	return args.Get(0).([]model.Story), args.Error(1)
}

func (mock *MockStoriesStore) UpdateStory(story *model.Story) (int64, error) {
	args := mock.Called(story)
	return args.Get(0).(int64), args.Error(1)
}

func (mock *MockStoriesStore) DeleteStory(storyID string) (int64, error) {
	args := mock.Called(storyID)
	return args.Get(0).(int64), args.Error(1)
}

func (mock *MockStoriesStore) GetMostViewsStories(offset, limit int) ([]model.Story, error) {
	args := mock.Called(offset, limit)
	return args.Get(0).([]model.Story), args.Error(1)
}

func (mock *MockStoriesStore) GetTopRatedStories(offset, limit int) ([]model.Story, error) {
	args := mock.Called(offset, limit)
	return args.Get(0).([]model.Story), args.Error(1)
}
