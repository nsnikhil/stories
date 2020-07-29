package store

import (
	"github.com/nsnikhil/stories/pkg/blog/domain"
	"go.uber.org/zap"
)

type TrieStoriesCache struct {
	trie   Trie
	logger *zap.Logger
}

func NewTrieStoriesCache(trie Trie, logger *zap.Logger) StoriesCache {
	return &TrieStoriesCache{
		trie:   trie,
		logger: logger,
	}
}

func (tc *TrieStoriesCache) AddStory(story *domain.Story) []error {
	res := make([]error, 0)

	addErr := tc.trie.insert(story.GetTitle(), story.GetID())
	for _, ar := range addErr {
		tc.logger.Debug(ar.Error(), zap.String("method", "AddStory"), zap.String("call", "insertTitle"))
		res = append(res, ar)
	}

	addErr = tc.trie.insert(story.GetBody(), story.GetID())
	for _, ar := range addErr {
		tc.logger.Debug(ar.Error(), zap.String("method", "AddStory"), zap.String("call", "insertBody"))
		res = append(res, ar)
	}

	return res
}

func (tc *TrieStoriesCache) GetStoryIDs(query string) ([]string, []error) {
	res := make([]error, 0)

	idm, getErr := tc.trie.getIDs(query)
	for _, ge := range getErr {
		tc.logger.Debug(ge.Error(), zap.String("method", "GetStoryIDs"))
		res = append(res, ge)
	}

	k := 0
	sz := len(idm)

	ids := make([]string, sz)
	for id := range idm {
		ids[k] = id
		k++
	}

	return ids, res
}
