package service

import "github.com/nsnikhil/stories/pkg/blog/dao"

type StoriesService interface {
	AddStory(story *dao.Story) error
	GetStory(storyID string) (*dao.Story, error)
	SearchStories(query string) ([]dao.Story, error)
	GetMostViewsStories(offset, limit int) ([]dao.Story, error)
	GetTopRatedStories(offset, limit int) ([]dao.Story, error)
}
