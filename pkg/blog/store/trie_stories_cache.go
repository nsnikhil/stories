package store

import (
	"github.com/nsnikhil/stories/cmd/config"
	"github.com/nsnikhil/stories/pkg/blog/domain"
	"go.uber.org/zap"
)

type TrieStoriesCache struct {
	trie   Trie
	config config.BlogConfig
	logger *zap.Logger
}

func NewTrieStoriesCache(trie Trie, config config.BlogConfig, logger *zap.Logger) StoriesCache {
	return &TrieStoriesCache{
		trie:   trie,
		config: config,
		logger: logger,
	}
}

func (tc *TrieStoriesCache) AddStory(story *domain.Story) []error {
	res := make([]error, 0)

	if tc.config.CacheTitle() {
		res = append(res, addToCache(story.GetTitle(), story.GetID(), "insertTitle", tc.logger, tc.trie)...)
	}

	if tc.config.CacheBody() {
		res = append(res, addToCache(story.GetBody(), story.GetID(), "insertBody", tc.logger, tc.trie)...)
	}

	return res
}

func addToCache(words, id, call string, lgr *zap.Logger, t Trie) []error {
	res := make([]error, 0)
	addErr := t.insert(words, id)
	for _, ar := range addErr {
		lgr.Debug(ar.Error(), zap.String("method", "AddStory"), zap.String("call", call))
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
