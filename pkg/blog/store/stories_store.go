package store

import (
	"github.com/nsnikhil/stories/pkg/blog/domain"
)

type StoriesStore interface {
	AddStory(story *domain.Story) error
	GetStories(storyIDs ...string) ([]domain.Story, error)
	UpdateStory(story *domain.Story) (int64, error)
	GetMostViewsStories(offset, limit int) ([]domain.Story, error)
	GetTopRatedStories(offset, limit int) ([]domain.Story, error)
}
