package story_test

import (
	"errors"
	"github.com/nsnikhil/stories/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestStoryServiceAddStory(t *testing.T) {
	testCases := []struct {
		name          string
		actualResult  func() error
		expectedError error
	}{
		{
			name: "test insert story success",
			actualResult: func() error {
				str, err := NewVanillaStory("title", "test body")
				require.NoError(t, err)

				mst := &str.MockStoriesStore{}
				mst.On("AddStory", str).Return(nil)

				mct := &str.MockStoriesCache{}
				mct.On("AddStory", str).Return([]error{})

				st := str.NewStore(mst, mct)
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.AddStory(str)
			},
		},
		{
			name: "test insert story failure",
			actualResult: func() error {
				str, err := NewVanillaStory("title", "test body")
				require.NoError(t, err)

				mst := &str.MockStoriesStore{}
				mst.On("AddStory", str).Return(errors.New("failed to insert story"))

				mct := &str.MockStoriesCache{}
				mct.On("AddStory", str).Return([]error{})

				st := str.NewStore(mst, mct)
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
		actualResult   func() (*Story, error)
		expectedResult func() *Story
		expectedError  error
	}{
		{
			name: "test get story success",
			actualResult: func() (*Story, error) {
				str, err := NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &str.MockStoriesStore{}
				mst.On("GetStories", []string{"2eaa0697-2572-47f9-bcff-0bdf0c7c6432"}).Return([]Story{*str}, nil)

				st := str.NewStore(mst, &str.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.GetStory("2eaa0697-2572-47f9-bcff-0bdf0c7c6432")
			},
			expectedResult: func() *Story {
				str, err := NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				return str
			},
		},
		{
			name: "test get story failure",
			actualResult: func() (*Story, error) {
				mst := &store.MockStoriesStore{}
				mst.On("GetStories", []string{"2eaa0697-2572-47f9-bcff-0bdf0c7c6432"}).Return([]Story{}, errors.New("failed to get story"))

				st := store.NewStore(mst, &store.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.GetStory("2eaa0697-2572-47f9-bcff-0bdf0c7c6432")
			},
			expectedResult: func() *Story {
				return nil
			},
			expectedError: errors.New("failed to get story"),
		},
		{
			name: "test get story failure when empty ids slice is empty",
			actualResult: func() (*Story, error) {
				mst := &store.MockStoriesStore{}
				mst.On("GetStories", []string{"2eaa0697-2572-47f9-bcff-0bdf0c7c6432"}).Return([]Story{}, nil)

				st := store.NewStore(mst, &store.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.GetStory("2eaa0697-2572-47f9-bcff-0bdf0c7c6432")
			},
			expectedResult: func() *Story {
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
				str, err := NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &str.MockStoriesStore{}
				mst.On("UpdateStory", str).Return(int64(1), nil)

				st := str.NewStore(mst, &str.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.UpdateStory(str)
			},
			expectedResult: 1,
		},
		{
			name: "test update story failure",
			actualResult: func() (int64, error) {
				str, err := NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &str.MockStoriesStore{}
				mst.On("UpdateStory", str).Return(int64(0), errors.New("failed to update story"))

				st := str.NewStore(mst, &str.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.UpdateStory(str)
			},
			expectedResult: 0,
			expectedError:  errors.New("failed to update story"),
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

func TestStoryServiceDeleteStory(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (int64, error)
		expectedResult int64
		expectedError  error
	}{
		{
			name: "test delete story success",
			actualResult: func() (int64, error) {
				str, err := NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &str.MockStoriesStore{}
				mst.On("DeleteStory", str.GetID()).Return(int64(1), nil)

				st := str.NewStore(mst, &str.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.DeleteStory(str.GetID())
			},
			expectedResult: 1,
		},
		{
			name: "test delete story failure",
			actualResult: func() (int64, error) {
				str, err := NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &str.MockStoriesStore{}
				mst.On("DeleteStory", str.GetID()).Return(int64(0), errors.New("failed to delete story"))

				st := str.NewStore(mst, &str.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.DeleteStory(str.GetID())
			},
			expectedResult: 0,
			expectedError:  errors.New("failed to delete story"),
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
		actualResult   func() ([]Story, error)
		expectedResult func() []Story
		expectedError  error
	}{
		{
			name: "test search story success",
			actualResult: func() ([]Story, error) {
				str, err := NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &str.MockStoriesStore{}
				mst.On("GetStories", []string{"2eaa0697-2572-47f9-bcff-0bdf0c7c6432"}).Return([]Story{*str}, nil)

				mct := &str.MockStoriesCache{}
				mct.On("GetStoryIDs", "body").Return([]string{"2eaa0697-2572-47f9-bcff-0bdf0c7c6432"}, []error{})

				st := str.NewStore(mst, mct)
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.SearchStories("body")
			},
			expectedResult: func() []Story {
				str, err := NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				return []Story{*str}
			},
		},
		{
			name: "test search story failure",
			actualResult: func() ([]Story, error) {
				mct := &store.MockStoriesCache{}
				mct.On("GetStoryIDs", "body").Return([]string{}, []error{errors.New("failed to find story id")})

				mst := &store.MockStoriesStore{}
				mst.On("GetStories", []string{"2eaa0697-2572-47f9-bcff-0bdf0c7c6432"}).Return([]Story{}, errors.New("failed to get story"))

				st := store.NewStore(mst, mct)
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.SearchStories("body")
			},
			expectedResult: func() []Story {
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
		actualResult   func() ([]Story, error)
		expectedResult func() []Story
		expectedError  error
	}{
		{
			name: "test get most viewed story success",
			actualResult: func() ([]Story, error) {
				str, err := NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &str.MockStoriesStore{}
				mst.On("GetMostViewsStories", 0, 1).Return([]Story{*str}, nil)

				st := str.NewStore(mst, &str.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.GetMostViewsStories(0, 1)
			},
			expectedResult: func() []Story {
				str, err := NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				return []Story{*str}
			},
		},
		{
			name: "test get most viewed story failure",
			actualResult: func() ([]Story, error) {
				mst := &store.MockStoriesStore{}
				mst.On("GetMostViewsStories", 0, 1).Return([]Story{}, errors.New("failed to get most viewed story"))

				st := store.NewStore(mst, &store.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.GetMostViewsStories(0, 1)
			},
			expectedResult: func() []Story {
				return []Story{}
			},
			expectedError: errors.New("failed to get most viewed story"),
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
		actualResult   func() ([]Story, error)
		expectedResult func() []Story
		expectedError  error
	}{
		{
			name: "test get top rated story success",
			actualResult: func() ([]Story, error) {
				str, err := NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &str.MockStoriesStore{}
				mst.On("GetTopRatedStories", 0, 1).Return([]Story{*str}, nil)

				st := str.NewStore(mst, &str.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.GetTopRatedStories(0, 1)
			},
			expectedResult: func() []Story {
				str, err := NewVanillaStory("title", "test body")
				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				return []Story{*str}
			},
		},
		{
			name: "test get top rated story failure",
			actualResult: func() ([]Story, error) {
				mst := &store.MockStoriesStore{}
				mst.On("GetTopRatedStories", 0, 1).Return([]Story{}, errors.New("failed to get top rated story"))

				st := store.NewStore(mst, &store.MockStoriesCache{})
				lgr := zap.NewExample()

				sv := NewDefaultStoriesService(st, lgr)

				return sv.GetTopRatedStories(0, 1)
			},
			expectedResult: func() []Story {
				return []Story{}
			},
			expectedError: errors.New("failed to get top rated story"),
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
