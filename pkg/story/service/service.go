package service

import (
	"errors"
	"fmt"
	"github.com/nsnikhil/stories/pkg/liberr"
	"github.com/nsnikhil/stories/pkg/store"
	"github.com/nsnikhil/stories/pkg/story/model"
)

type StoryService interface {
	AddStory(story *model.Story) error
	GetStory(storyID string) (*model.Story, error)

	UpdateStory(story *model.Story) (int64, error)
	DeleteStory(storyID string) (int64, error)

	SearchStories(query string) ([]model.Story, error)

	GetMostViewsStories(offset, limit int) ([]model.Story, error)
	GetTopRatedStories(offset, limit int) ([]model.Story, error)
}

//TODO: RENAME (REMOVE DEFAULT)
type defaultStoriesService struct {
	store store.StoriesStore
}

//TODO: REMOVE ERROR NIL CHECK JUST TO INJECT OPERATIONS IN THIS AND ALL THE METHODS BELOW
func (dss *defaultStoriesService) AddStory(story *model.Story) error {
	_, err := dss.store.AddStory(story)
	if err != nil {
		return liberr.WithOperation(op("AddStory"), err)
	}

	return nil
}

func (dss *defaultStoriesService) GetStory(storyID string) (*model.Story, error) {
	stories, err := dss.store.GetStories(storyID)
	if err != nil {
		return nil, liberr.WithOperation(op("GetStory"), err)
	}

	return &stories[0], nil
}

func (dss *defaultStoriesService) UpdateStory(story *model.Story) (int64, error) {
	c, err := dss.store.UpdateStory(story)
	if err != nil {
		return 0, liberr.WithOperation(op("UpdateStory"), err)
	}

	return c, err
}

func (dss *defaultStoriesService) DeleteStory(storyID string) (int64, error) {
	c, err := dss.store.DeleteStory(storyID)
	if err != nil {
		return 0, liberr.WithOperation(op("DeleteStory"), err)
	}

	return c, err
}

//TODO: FINISH THE IMPLEMENTATION
func (dss *defaultStoriesService) SearchStories(query string) ([]model.Story, error) {
	return nil, errors.New("UNIMPLEMENTED")
}

func (dss *defaultStoriesService) GetMostViewsStories(offset, limit int) ([]model.Story, error) {
	res, err := dss.store.GetMostViewsStories(offset, limit)
	if err != nil {
		return nil, liberr.WithOperation(op("GetMostViewsStories"), err)
	}

	return res, nil
}

func (dss *defaultStoriesService) GetTopRatedStories(offset, limit int) ([]model.Story, error) {
	res, err := dss.store.GetTopRatedStories(offset, limit)
	if err != nil {
		return nil, liberr.WithOperation(op("GetTopRatedStories"), err)
	}

	return res, nil
}

//TODO: REMOVE THIS HELPER FUNCTION
func op(co string) liberr.Operation {
	return liberr.Operation(fmt.Sprintf("StoryService.%s", co))
}

func NewStoriesService(store store.StoriesStore) StoryService {
	return &defaultStoriesService{
		store: store,
	}
}
