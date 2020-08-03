package store

import (
	"errors"
	"github.com/nsnikhil/stories/cmd/config"
	"github.com/nsnikhil/stories/pkg/blog/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestCreateNewStoriesCache(t *testing.T) {
	tr := &mockTrie{}
	lgr := zap.NewExample()
	cfg := config.LoadConfigs().GetBlogConfig()

	actualResult := NewTrieStoriesCache(tr, cfg, lgr)
	expectedResult := &TrieStoriesCache{trie: tr, config: cfg, logger: lgr}

	assert.Equal(t, expectedResult, actualResult)
}

func TestStoriesCacheAddStory(t *testing.T) {
	cfg := config.LoadConfigs().GetBlogConfig()

	testCases := []struct {
		name          string
		actualResult  func() []error
		expectedError []error
	}{
		{
			name: "test insert into cache",
			actualResult: func() []error {
				tr := &mockTrie{}
				tr.On("insert", "title", "some-id").Return([]error{})
				tr.On("insert", "test body", "some-id").Return([]error{})
				tc := NewTrieStoriesCache(tr, cfg, zap.NewExample())

				st, err := domain.NewVanillaStory("title", "test body")
				require.NoError(t, err)
				st.ID = "some-id"

				return tc.AddStory(st)
			},
			expectedError: []error{},
		},
		{
			name: "test insert into cache return errors when title contains symbols",
			actualResult: func() []error {
				tr := &mockTrie{}
				tr.On("insert", "title @", "some-id").Return([]error{errors.New("@ is not a valid character")})
				tr.On("insert", "test body", "some-id").Return([]error{})
				tc := NewTrieStoriesCache(tr, cfg, zap.NewExample())

				st, err := domain.NewVanillaStory("title @", "test body")
				require.NoError(t, err)
				st.ID = "some-id"

				return tc.AddStory(st)
			},
			expectedError: []error{errors.New("@ is not a valid character")},
		},
		{
			name: "test insert into cache return errors when body contains symbols",
			actualResult: func() []error {
				tr := &mockTrie{}
				tr.On("insert", "title", "some-id").Return([]error{})
				tr.On("insert", "test @ some value", "some-id").Return([]error{errors.New("@ is not a valid character")})
				tc := NewTrieStoriesCache(tr, cfg, zap.NewExample())

				st, err := domain.NewVanillaStory("title", "test @ some value")
				require.NoError(t, err)
				st.ID = "some-id"

				return tc.AddStory(st)
			},
			expectedError: []error{errors.New("@ is not a valid character")},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedError, testCase.actualResult())
		})
	}
}

func TestStoriesCacheGetStoryIDs(t *testing.T) {
	cfg := config.LoadConfigs().GetBlogConfig()

	testCases := []struct {
		name           string
		actualResult   func() ([]string, []error)
		expectedResult []string
		expectedError  []error
	}{
		{
			name: "test get from cache",
			actualResult: func() ([]string, []error) {
				tr := &mockTrie{}
				tr.On("getIDs", "test").Return(map[string]bool{"36982b87-be33-4683-aaaa-e69282a03c83": true}, []error{})
				tc := NewTrieStoriesCache(tr, cfg, zap.NewExample())

				return tc.GetStoryIDs("test")
			},
			expectedResult: []string{"36982b87-be33-4683-aaaa-e69282a03c83"},
			expectedError:  []error{},
		},
		{
			name: "test get from cache return error",
			actualResult: func() ([]string, []error) {
				tr := &mockTrie{}
				tr.On("getIDs", "search this @").Return(map[string]bool{"36982b87-be33-4683-aaaa-e69282a03c83": true}, []error{errors.New("@ is not a valid character")})
				tc := NewTrieStoriesCache(tr, cfg, zap.NewExample())

				return tc.GetStoryIDs("search this @")
			},
			expectedResult: []string{"36982b87-be33-4683-aaaa-e69282a03c83"},
			expectedError:  []error{errors.New("@ is not a valid character")},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult, res)
		})
	}
}
