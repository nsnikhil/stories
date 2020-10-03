package service

import (
	"errors"
	"github.com/nsnikhil/stories/pkg/store"
	"github.com/nsnikhil/stories/pkg/story/model"
)

//goland:noinspection ALL
type StoryService interface {
	AddStory(story *model.Story) error
	GetStory(storyID string) (*model.Story, error)

	UpdateStory(story *model.Story) (int64, error)
	DeleteStory(storyID string) (int64, error)

	SearchStories(query string) ([]model.Story, error)

	GetMostViewsStories(offset, limit int) ([]model.Story, error)
	GetTopRatedStories(offset, limit int) ([]model.Story, error)
}

type defaultStoriesService struct {
	store store.StoriesStore
}

func (dss *defaultStoriesService) AddStory(story *model.Story) error {
	_, err := dss.store.AddStory(story)
	if err != nil {
		return err
	}

	return nil
}

func (dss *defaultStoriesService) GetStory(storyID string) (*model.Story, error) {
	stories, err := dss.store.GetStories(storyID)
	if err != nil {
		return nil, err
	}

	return &stories[0], nil
}

func (dss *defaultStoriesService) UpdateStory(story *model.Story) (int64, error) {
	return dss.store.UpdateStory(story)
}

func (dss *defaultStoriesService) DeleteStory(storyID string) (int64, error) {
	return dss.store.DeleteStory(storyID)
}

func (dss *defaultStoriesService) SearchStories(query string) ([]model.Story, error) {
	return nil, errors.New("UNIMPLEMENTED")
}

func (dss *defaultStoriesService) GetMostViewsStories(offset, limit int) ([]model.Story, error) {
	return dss.store.GetMostViewsStories(offset, limit)
}

func (dss *defaultStoriesService) GetTopRatedStories(offset, limit int) ([]model.Story, error) {
	return dss.store.GetTopRatedStories(offset, limit)
}

func NewStoriesService(store store.StoriesStore) StoryService {
	return &defaultStoriesService{
		store: store,
	}
}
