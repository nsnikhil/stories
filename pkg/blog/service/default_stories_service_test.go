package service

import (
	"errors"
	"github.com/nsnikhil/stories/pkg/blog/domain"
	"github.com/nsnikhil/stories/pkg/blog/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestCreateNewStoriesService(t *testing.T) {
	st := store.NewStore(&store.MockStoriesStore{}, &store.MockStoriesCache{})
	lgr := zap.NewExample()

	actualResult := NewDefaultStoriesService(st, lgr)
	expectedResult := &DefaultStoriesService{store: st, logger: lgr}

	assert.Equal(t, expectedResult, actualResult)
}

func TestStoryServiceAddStory(t *testing.T) {
	testCases := []struct {
		name          string
		actualResult  func() error
		expectedError error
	}{
		{
			name: "test insert story success",
			actualResult: func() error {
				str, err := domain.NewVanillaStory("title", "test body")
				require.NoError(t, err)

				mst := &store.MockStoriesStore{}
				mst.On("AddStory", str).Return(nil)

				mct := &store.MockStoriesCache{}
				mct.On("AddStory", str).Return([]error{})

				st := store.NewStore(mst, mct)
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.AddStory(str)
			},
		},
		{
			name: "test insert story failure",
			actualResult: func() error {
				str, err := domain.NewVanillaStory("title", "test body")
				require.NoError(t, err)

				mst := &store.MockStoriesStore{}
				mst.On("AddStory", str).Return(errors.New("failed to insert story"))

				mct := &store.MockStoriesCache{}
				mct.On("AddStory", str).Return([]error{})

				st := store.NewStore(mst, mct)
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.AddStory(str)
			},
			expectedError: errors.New("failed to insert story"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedError, testCase.actualResult())
		})
	}
}

func TestStoryServiceGetStory(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (*domain.Story, error)
		expectedResult func() *domain.Story
		expectedError  error
	}{
		{
			name: "test get story success",
			actualResult: func() (*domain.Story, error) {
				str, err := domain.NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.Id = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("GetStories", []string{"2eaa0697-2572-47f9-bcff-0bdf0c7c6432"}).Return([]domain.Story{*str}, nil)

				st := store.NewStore(mst, &store.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.GetStory("2eaa0697-2572-47f9-bcff-0bdf0c7c6432")
			},
			expectedResult: func() *domain.Story {
				str, err := domain.NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.Id = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				return str
			},
		},
		{
			name: "test get story failure",
			actualResult: func() (*domain.Story, error) {
				mst := &store.MockStoriesStore{}
				mst.On("GetStories", []string{"2eaa0697-2572-47f9-bcff-0bdf0c7c6432"}).Return([]domain.Story{}, errors.New("failed to get story"))

				st := store.NewStore(mst, &store.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.GetStory("2eaa0697-2572-47f9-bcff-0bdf0c7c6432")
			},
			expectedResult: func() *domain.Story {
				return nil
			},
			expectedError: errors.New("failed to get story"),
		},
		{
			name: "test get story failure when empty ids slice is empty",
			actualResult: func() (*domain.Story, error) {
				mst := &store.MockStoriesStore{}
				mst.On("GetStories", []string{"2eaa0697-2572-47f9-bcff-0bdf0c7c6432"}).Return([]domain.Story{}, nil)

				st := store.NewStore(mst, &store.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.GetStory("2eaa0697-2572-47f9-bcff-0bdf0c7c6432")
			},
			expectedResult: func() *domain.Story {
				return nil
			},
			expectedError: errors.New("no story found for id 2eaa0697-2572-47f9-bcff-0bdf0c7c6432"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult(), res)
		})
	}
}

func TestStoryServiceUpdateStory(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (int64, error)
		expectedResult int64
		expectedError  error
	}{
		{
			name: "test update story success",
			actualResult: func() (int64, error) {
				str, err := domain.NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.Id = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("UpdateStory", str).Return(int64(1), nil)

				st := store.NewStore(mst, &store.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.UpdateStory(str)
			},
			expectedResult: 1,
		},
		{
			name: "test update story failure",
			actualResult: func() (int64, error) {
				str, err := domain.NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.Id = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("UpdateStory", str).Return(int64(0), errors.New("failed to update story"))

				st := store.NewStore(mst, &store.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.UpdateStory(str)
			},
			expectedResult: 0,
			expectedError: errors.New("failed to update story"),
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

func TestStoryServiceSearchStories(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() ([]domain.Story, error)
		expectedResult func() []domain.Story
		expectedError  error
	}{
		{
			name: "test search stories success",
			actualResult: func() ([]domain.Story, error) {
				str, err := domain.NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.Id = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("GetStories", []string{"2eaa0697-2572-47f9-bcff-0bdf0c7c6432"}).Return([]domain.Story{*str}, nil)

				mct := &store.MockStoriesCache{}
				mct.On("GetStoryIDs", "body").Return([]string{"2eaa0697-2572-47f9-bcff-0bdf0c7c6432"}, []error{})

				st := store.NewStore(mst, mct)
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.SearchStories("body")
			},
			expectedResult: func() []domain.Story {
				str, err := domain.NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.Id = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				return []domain.Story{*str}
			},
		},
		{
			name: "test search stories failure",
			actualResult: func() ([]domain.Story, error) {
				mct := &store.MockStoriesCache{}
				mct.On("GetStoryIDs", "body").Return([]string{}, []error{errors.New("failed to find story id")})

				mst := &store.MockStoriesStore{}
				mst.On("GetStories", []string{"2eaa0697-2572-47f9-bcff-0bdf0c7c6432"}).Return([]domain.Story{}, errors.New("failed to get story"))

				st := store.NewStore(mst, mct)
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.SearchStories("body")
			},
			expectedResult: func() []domain.Story {
				return nil
			},
			expectedError: errors.New("failed to find story id"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult(), res)
		})
	}
}

func TestStoryServiceGetMostViewsStories(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() ([]domain.Story, error)
		expectedResult func() []domain.Story
		expectedError  error
	}{
		{
			name: "test get most viewed stories success",
			actualResult: func() ([]domain.Story, error) {
				str, err := domain.NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.Id = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("GetMostViewsStories", 0, 1).Return([]domain.Story{*str}, nil)

				st := store.NewStore(mst, &store.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.GetMostViewsStories(0, 1)
			},
			expectedResult: func() []domain.Story {
				str, err := domain.NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.Id = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				return []domain.Story{*str}
			},
		},
		{
			name: "test get most viewed stories failure",
			actualResult: func() ([]domain.Story, error) {
				mst := &store.MockStoriesStore{}
				mst.On("GetMostViewsStories", 0, 1).Return([]domain.Story{}, errors.New("failed to get most viewed stories"))

				st := store.NewStore(mst, &store.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.GetMostViewsStories(0, 1)
			},
			expectedResult: func() []domain.Story {
				return []domain.Story{}
			},
			expectedError: errors.New("failed to get most viewed stories"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult(), res)
		})
	}
}

func TestStoryServiceGetTopRatedStories(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() ([]domain.Story, error)
		expectedResult func() []domain.Story
		expectedError  error
	}{
		{
			name: "test get top rated stories success",
			actualResult: func() ([]domain.Story, error) {
				str, err := domain.NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.Id = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("GetTopRatedStories", 0, 1).Return([]domain.Story{*str}, nil)

				st := store.NewStore(mst, &store.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.GetTopRatedStories(0, 1)
			},
			expectedResult: func() []domain.Story {
				str, err := domain.NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.Id = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				return []domain.Story{*str}
			},
		},
		{
			name: "test get top rated stories failure",
			actualResult: func() ([]domain.Story, error) {
				mst := &store.MockStoriesStore{}
				mst.On("GetTopRatedStories", 0, 1).Return([]domain.Story{}, errors.New("failed to get top rated stories"))

				st := store.NewStore(mst, &store.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.GetTopRatedStories(0, 1)
			},
			expectedResult: func() []domain.Story {
				return []domain.Story{}
			},
			expectedError: errors.New("failed to get top rated stories"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult(), res)
		})
	}
}
