package service

import "github.com/nsnikhil/stories/pkg/blog/domain"

type StoriesService interface {
	AddStory(story *domain.Story) error
	GetStory(storyID string) (*domain.Story, error)
	SearchStories(query string) ([]domain.Story, error)
	GetMostViewsStories(offset, limit int) ([]domain.Story, error)
	GetTopRatedStories(offset, limit int) ([]domain.Story, error)
}
