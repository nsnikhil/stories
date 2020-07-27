package store

import "github.com/nsnikhil/stories/pkg/blog/dto"

type StoriesStore interface {
	AddStory(story *dto.Story) error
	GetStories(storyIDs ...string) ([]dto.Story, error)
	GetMostViewsStories(offset, limit int) ([]dto.Story, error)
	GetTopRatedStories(offset, limit int) ([]dto.Story, error)
}

type StoriesCache interface {
	AddStory(story *dto.Story) error
	GetStoryIDs(query string) ([]string, error)
}
