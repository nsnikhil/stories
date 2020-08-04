package store

import (
	"github.com/nsnikhil/stories/pkg/blog/domain"
	"github.com/stretchr/testify/mock"
)

type MockStoriesStore struct {
	mock.Mock
}

func (mock *MockStoriesStore) AddStory(story *domain.Story) error {
	args := mock.Called(story)
	return args.Error(0)
}

func (mock *MockStoriesStore) GetStories(storyIDs ...string) ([]domain.Story, error) {
	args := mock.Called(storyIDs)
	return args.Get(0).([]domain.Story), args.Error(1)
}

func (mock *MockStoriesStore) UpdateStory(story *domain.Story) (int64, error) {
	args := mock.Called(story)
	return args.Get(0).(int64), args.Error(1)
}

func (mock *MockStoriesStore) DeleteStory(storyID string) (int64, error) {
	args := mock.Called(storyID)
	return args.Get(0).(int64), args.Error(1)
}

func (mock *MockStoriesStore) GetMostViewsStories(offset, limit int) ([]domain.Story, error) {
	args := mock.Called(offset, limit)
	return args.Get(0).([]domain.Story), args.Error(1)
}

func (mock *MockStoriesStore) GetTopRatedStories(offset, limit int) ([]domain.Story, error) {
	args := mock.Called(offset, limit)
	return args.Get(0).([]domain.Story), args.Error(1)
}

type MockStoriesCache struct {
	mock.Mock
}

func (mock *MockStoriesCache) AddStory(story *domain.Story) []error {
	args := mock.Called(story)
	return args.Get(0).([]error)
}

func (mock *MockStoriesCache) GetStoryIDs(query string) ([]string, []error) {
	args := mock.Called(query)
	return args.Get(0).([]string), args.Get(1).([]error)
}

type mockTrie struct {
	mock.Mock
}

func (mock *mockTrie) insert(s, id string) []error {
	args := mock.Called(s, id)
	return args.Get(0).([]error)
}

func (mock *mockTrie) getIDs(query string) (map[string]bool, []error) {
	args := mock.Called(query)
	return args.Get(0).(map[string]bool), args.Get(1).([]error)
}
