package service

import (
	"fmt"
	"github.com/nsnikhil/stories/pkg/blog/domain"
	"github.com/nsnikhil/stories/pkg/blog/store"
	"go.uber.org/zap"
)

type DefaultStoriesService struct {
	store  *store.Store
	logger *zap.Logger
}

func NewDefaultStoriesService(store *store.Store, logger *zap.Logger) StoriesService {
	return &DefaultStoriesService{
		store:  store,
		logger: logger,
	}
}

func (dss *DefaultStoriesService) AddStory(story *domain.Story) error {
	err := dss.store.GetStoriesStore().AddStory(story)
	if err != nil {
		return err
	}

	go dss.store.GetStoriesCache().AddStory(story)

	return nil

}

func (dss *DefaultStoriesService) GetStory(storyID string) (*domain.Story, error) {
	stories, err := dss.store.GetStoriesStore().GetStories(storyID)
	if err != nil {
		return nil, err
	}

	if len(stories) == 0 {
		err := fmt.Errorf("no story found for id %s", storyID)
		dss.logger.Error(err.Error(), zap.String("method", "GetStory"))
		return nil, err
	}

	return &stories[0], nil
}

func (dss *DefaultStoriesService) UpdateStory(story *domain.Story) (int64, error) {
	return dss.store.GetStoriesStore().UpdateStory(story)
}

func (dss *DefaultStoriesService) SearchStories(query string) ([]domain.Story, error) {
	ids, errs := dss.store.GetStoriesCache().GetStoryIDs(query)
	if len(ids) == 0 && len(errs) != 0 {
		return nil, errs[0]
	}

	return dss.store.GetStoriesStore().GetStories(ids...)
}

func (dss *DefaultStoriesService) GetMostViewsStories(offset, limit int) ([]domain.Story, error) {
	return dss.store.GetStoriesStore().GetMostViewsStories(offset, limit)
}

func (dss *DefaultStoriesService) GetTopRatedStories(offset, limit int) ([]domain.Story, error) {
	return dss.store.GetStoriesStore().GetTopRatedStories(offset, limit)
}
