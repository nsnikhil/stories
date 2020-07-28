package store

import (
	"github.com/nsnikhil/stories/pkg/blog/domain"
)

type StoriesCache interface {
	AddStory(story *domain.Story) []error
	GetStoryIDs(query string) ([]string, []error)
}
