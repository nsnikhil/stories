package service_test

import (
	"errors"
	"github.com/nsnikhil/stories/pkg/liberr"
	"github.com/nsnikhil/stories/pkg/store"
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/nsnikhil/stories/pkg/story/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStoryServiceAddStory(t *testing.T) {
	testCases := map[string]struct {
		store         func() store.StoriesStore
		expectedError error
	}{
		"test add story success": {
			store: func() store.StoriesStore {
				mst := &store.MockStoriesStore{}
				mst.On("AddStory", mock.AnythingOfType("*model.Story")).Return("a45c9dac-56dc-4771-a3f4-f10ad30a20a5", nil)

				return mst
			},
			expectedError: nil,
		},
		"test add story failure when dependency fails": {
			store: func() store.StoriesStore {
				mst := &store.MockStoriesStore{}
				mst.On("AddStory", mock.AnythingOfType("*model.Story")).Return("", liberr.WithArgs(liberr.SeverityError, errors.New("failed to insert story")))

				return mst
			},
			expectedError: errors.New("failed to insert story"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			svc := service.NewStoriesService(testCase.store())

			str, err := model.NewStoryBuilder().
				SetTitle(100, "title").
				SetBody(100, "test body").
				Build()

			require.NoError(t, err)

			err = svc.AddStory(str)

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestStoryServiceGetStory(t *testing.T) {
	testCases := map[string]struct {
		input         func() (store.StoriesStore, string)
		expectedStory func() *model.Story
		expectedError error
	}{
		"test get story success": {
			input: func() (store.StoriesStore, string) {
				id := "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = id

				mst := &store.MockStoriesStore{}
				mst.On("GetStories", []string{id}).Return([]model.Story{*str}, nil)

				return mst, id
			},
			expectedStory: func() *model.Story {
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				return str
			},
		},
		"test get story failure": {
			input: func() (store.StoriesStore, string) {
				id := "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("GetStories", []string{id}).Return([]model.Story{}, liberr.WithArgs(errors.New("failed to get story")))

				return mst, id
			},
			expectedStory: func() *model.Story {
				return nil
			},
			expectedError: errors.New("failed to get story"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			st, id := testCase.input()

			svc := service.NewStoriesService(st)

			res, err := svc.GetStory(id)

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, testCase.expectedStory(), res)
		})
	}
}

func TestStoryServiceUpdateStory(t *testing.T) {
	testCases := map[string]struct {
		input         func() (*model.Story, store.StoriesStore)
		expectedCount int64
		expectedError error
	}{
		"test update story success": {
			input: func() (*model.Story, store.StoriesStore) {
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("UpdateStory", str).Return(int64(1), nil)

				return str, mst
			},
			expectedCount: 1,
		},
		"test update story failure": {
			input: func() (*model.Story, store.StoriesStore) {
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("UpdateStory", str).Return(int64(0), liberr.WithArgs(errors.New("failed to update story")))

				return str, mst
			},
			expectedCount: 0,
			expectedError: errors.New("failed to update story"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			st, str := testCase.input()

			svc := service.NewStoriesService(str)

			res, err := svc.UpdateStory(st)

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, testCase.expectedCount, res)
		})
	}
}

func TestStoryServiceDeleteStory(t *testing.T) {
	testCases := map[string]struct {
		input         func() (string, store.StoriesStore)
		expectedCount int64
		expectedError error
	}{
		"test delete story success": {
			input: func() (string, store.StoriesStore) {
				id := "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = id

				mst := &store.MockStoriesStore{}
				mst.On("DeleteStory", str.GetID()).Return(int64(1), nil)

				return id, mst
			},
			expectedCount: 1,
		},
		"test delete story failure": {
			input: func() (string, store.StoriesStore) {
				id := "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = id

				mst := &store.MockStoriesStore{}
				mst.On("DeleteStory", str.GetID()).Return(int64(0), liberr.WithArgs(errors.New("failed to delete story")))

				return id, mst
			},
			expectedCount: 0,
			expectedError: errors.New("failed to delete story"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			id, str := testCase.input()

			svc := service.NewStoriesService(str)

			res, err := svc.DeleteStory(id)

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, testCase.expectedCount, res)
		})
	}
}

func TestStoryServiceSearchStories(t *testing.T) {
	// TODO
}

func TestStoryServiceGetMostViewsStories(t *testing.T) {
	testCases := map[string]struct {
		input          func() (int, int, store.StoriesStore)
		expectedResult func() []model.Story
		expectedError  error
	}{
		"test get most viewed story success": {
			input: func() (int, int, store.StoriesStore) {
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("GetMostViewsStories", 0, 1).Return([]model.Story{*str}, nil)

				return 0, 1, mst
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
		"test get most viewed story failure": {
			input: func() (int, int, store.StoriesStore) {
				mst := &store.MockStoriesStore{}
				mst.On("GetMostViewsStories", 0, 1).Return([]model.Story{}, liberr.WithArgs(liberr.SeverityError, errors.New("failed to get most viewed story")))

				return 0, 1, mst
			},
			expectedResult: func() []model.Story {
				return nil
			},
			expectedError: errors.New("failed to get most viewed story"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			o, l, st := testCase.input()

			svc := service.NewStoriesService(st)

			res, err := svc.GetMostViewsStories(o, l)

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
	testCases := map[string]struct {
		input          func() (int, int, store.StoriesStore)
		expectedResult func() []model.Story
		expectedError  error
	}{
		"test get top rated story success": {
			input: func() (int, int, store.StoriesStore) {
				str, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)
				str.ID = "2eaa0697-2572-47f9-bcff-0bdf0c7c6432"

				mst := &store.MockStoriesStore{}
				mst.On("GetTopRatedStories", 0, 1).Return([]model.Story{*str}, nil)

				return 0, 1, mst
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
		"test get top rated story failure": {
			input: func() (int, int, store.StoriesStore) {
				mst := &store.MockStoriesStore{}
				mst.On("GetTopRatedStories", 0, 1).Return([]model.Story{}, liberr.WithArgs(liberr.SeverityError, errors.New("failed to get top rated story")))

				return 0, 1, mst
			},
			expectedResult: func() []model.Story {
				return nil
			},
			expectedError: errors.New("failed to get top rated story"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			o, l, st := testCase.input()

			svc := service.NewStoriesService(st)

			res, err := svc.GetTopRatedStories(o, l)

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, testCase.expectedResult(), res)
		})
	}
}
