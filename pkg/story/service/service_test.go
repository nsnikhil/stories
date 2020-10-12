package service_test

import (
	"errors"
	"github.com/nsnikhil/stories/pkg/liberr"
	"github.com/nsnikhil/stories/pkg/store"
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/nsnikhil/stories/pkg/story/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)

				mst := &store.MockStoriesStore{}
				mst.On("AddStory", str).Return("a45c9dac-56dc-4771-a3f4-f10ad30a20a5", nil)

				sv := service.NewStoriesService(mst)

				return sv.AddStory(str)
			},
		},
		{
			name: "test insert story failure",
			actualResult: func() error {
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)

				mst := &store.MockStoriesStore{}
				mst.On("AddStory", str).Return("", liberr.WithArgs(liberr.SeverityError, errors.New("failed to insert story")))

				sv := service.NewStoriesService(mst)

				return sv.AddStory(str)
			},
			expectedError: errors.New("failed to insert story"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.actualResult()
			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestStoryServiceGetStory(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (*model.Story, error)
		expectedResult func() *model.Story
		expectedError  error
	}{
		{
			name: "test get story success",
			actualResult: func() (*model.Story, error) {
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("GetStories", []string{"2eaa0697-2572-47f9-bcff-0bdf0c7c6432"}).Return([]model.Story{*str}, nil)

				sv := service.NewStoriesService(mst)

				return sv.GetStory("2eaa0697-2572-47f9-bcff-0bdf0c7c6432")
			},
			expectedResult: func() *model.Story {
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				return str
			},
		},
		{
			name: "test get story failure",
			actualResult: func() (*model.Story, error) {
				mst := &store.MockStoriesStore{}
				mst.On("GetStories", []string{"2eaa0697-2572-47f9-bcff-0bdf0c7c6432"}).Return([]model.Story{}, liberr.WithArgs(liberr.SeverityError, errors.New("failed to get story")))

				sv := service.NewStoriesService(mst)

				return sv.GetStory("2eaa0697-2572-47f9-bcff-0bdf0c7c6432")
			},
			expectedResult: func() *model.Story {
				return nil
			},
			expectedError: errors.New("failed to get story"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}

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
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("UpdateStory", str).Return(int64(1), nil)

				sv := service.NewStoriesService(mst)

				return sv.UpdateStory(str)
			},
			expectedResult: 1,
		},
		{
			name: "test update story failure",
			actualResult: func() (int64, error) {
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("UpdateStory", str).Return(int64(0), liberr.WithArgs(liberr.SeverityError, errors.New("failed to update story")))

				sv := service.NewStoriesService(mst)

				return sv.UpdateStory(str)
			},
			expectedResult: 0,
			expectedError:  errors.New("failed to update story"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}

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
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("DeleteStory", str.GetID()).Return(int64(1), nil)

				sv := service.NewStoriesService(mst)

				return sv.DeleteStory(str.GetID())
			},
			expectedResult: 1,
		},
		{
			name: "test delete story failure",
			actualResult: func() (int64, error) {
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("DeleteStory", str.GetID()).Return(int64(0), liberr.WithArgs(liberr.SeverityError, errors.New("failed to delete story")))

				sv := service.NewStoriesService(mst)

				return sv.DeleteStory(str.GetID())
			},
			expectedResult: 0,
			expectedError:  errors.New("failed to delete story"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, testCase.expectedResult, res)
		})
	}
}

func TestStoryServiceSearchStories(t *testing.T) {
	// TODO
}

func TestStoryServiceGetMostViewsStories(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() ([]model.Story, error)
		expectedResult func() []model.Story
		expectedError  error
	}{
		{
			name: "test get most viewed story success",
			actualResult: func() ([]model.Story, error) {
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("GetMostViewsStories", 0, 1).Return([]model.Story{*str}, nil)

				sv := service.NewStoriesService(mst)

				return sv.GetMostViewsStories(0, 1)
			},
			expectedResult: func() []model.Story {
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				return []model.Story{*str}
			},
		},
		{
			name: "test get most viewed story failure",
			actualResult: func() ([]model.Story, error) {
				mst := &store.MockStoriesStore{}
				mst.On("GetMostViewsStories", 0, 1).Return([]model.Story{}, liberr.WithArgs(liberr.SeverityError, errors.New("failed to get most viewed story")))

				sv := service.NewStoriesService(mst)

				return sv.GetMostViewsStories(0, 1)
			},
			expectedResult: func() []model.Story {
				return nil
			},
			expectedError: errors.New("failed to get most viewed story"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, testCase.expectedResult(), res)
		})
	}
}

func TestStoryServiceGetTopRatedStories(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() ([]model.Story, error)
		expectedResult func() []model.Story
		expectedError  error
	}{
		{
			name: "test get top rated story success",
			actualResult: func() ([]model.Story, error) {
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("GetTopRatedStories", 0, 1).Return([]model.Story{*str}, nil)

				sv := service.NewStoriesService(mst)

				return sv.GetTopRatedStories(0, 1)
			},
			expectedResult: func() []model.Story {
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				return []model.Story{*str}
			},
		},
		{
			name: "test get top rated story failure",
			actualResult: func() ([]model.Story, error) {
				mst := &store.MockStoriesStore{}
				mst.On("GetTopRatedStories", 0, 1).Return([]model.Story{}, liberr.WithArgs(liberr.SeverityError, errors.New("failed to get top rated story")))

				sv := service.NewStoriesService(mst)

				return sv.GetTopRatedStories(0, 1)
			},
			expectedResult: func() []model.Story {
				return nil
			},
			expectedError: errors.New("failed to get top rated story"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, testCase.expectedResult(), res)
		})
	}
}
