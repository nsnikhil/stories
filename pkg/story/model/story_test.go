package model_test

import (
	"errors"
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateNewStory(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (*model.Story, error)
		expectedResult *model.Story
		expectedError  error
	}{
		{
			name: "test create new story with title and body",
			actualResult: func() (*model.Story, error) {
				return model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(10000, "this is a test body").
					Build()
			},
			expectedResult: &model.Story{
				Title: "title",
				Body:  "this is a test body",
			},
		},
		{
			name: "test create new story with title and body and id",
			actualResult: func() (*model.Story, error) {
				return model.NewStoryBuilder().
					SetID("ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a").
					SetTitle(100, "title").
					SetBody(10000, "this is a test body").
					Build()
			},
			expectedResult: &model.Story{
				ID:    "ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a",
				Title: "title",
				Body:  "this is a test body",
			},
		},
		{
			name: "test failed to create story when id is invalid",
			actualResult: func() (*model.Story, error) {
				return model.NewStoryBuilder().
					SetID("invalid").
					SetTitle(100, "title").
					SetBody(10000, "this is a test body").
					Build()
			},
			expectedError: errors.New("invalid id: invalid"),
		},
		{
			name: "test failed to create story when title is empty",
			actualResult: func() (*model.Story, error) {
				return model.NewStoryBuilder().
					SetTitle(100, "").
					SetBody(10000, "this is a test body").
					Build()
			},
			expectedError: errors.New("title cannot be empty"),
		},
		{
			name: "test failed to create story when title exceeds max length",
			actualResult: func() (*model.Story, error) {
				return model.NewStoryBuilder().
					SetTitle(10, "this is a very long title").
					SetBody(10000, "this is a test body").
					Build()
			},
			expectedError: errors.New("title max length exceeded"),
		},
		{
			name: "test failed to create story when title is empty",
			actualResult: func() (*model.Story, error) {
				return model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(10000, "").
					Build()
			},
			expectedError: errors.New("body cannot be empty"),
		},
		{
			name: "test failed to create story when title exceeds max length",
			actualResult: func() (*model.Story, error) {
				return model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(10, "this is a test body").
					Build()
			},
			expectedError: errors.New("body max length exceeded"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			}else {
				assert.Nil(t, err)
			}

			assert.Equal(t, testCase.expectedResult, res)
		})
	}
}

func TestStoryGetter(t *testing.T) {
	id := "ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a"
	createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

	st, err := model.NewStoryBuilder().
		SetID(id).
		SetTitle(100, "title").
		SetBody(10000, "this is a test body").
		SetViewCount(25).
		SetUpVotes(10).
		SetDownVotes(2).
		SetCreatedAt(createdAt).
		SetUpdatedAt(updatedAt).
		Build()

	require.NoError(t, err)

	testCases := []struct {
		name           string
		actualResult   interface{}
		expectedResult interface{}
	}{
		{
			name:           "test get id",
			actualResult:   st.GetID(),
			expectedResult: "ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a",
		},
		{
			name:           "test get title",
			actualResult:   st.GetTitle(),
			expectedResult: "title",
		},
		{
			name:           "test get body",
			actualResult:   st.GetBody(),
			expectedResult: "this is a test body",
		},
		{
			name:           "test get view count",
			actualResult:   st.GetViewCount(),
			expectedResult: int64(25),
		},
		{
			name:           "test get up votes",
			actualResult:   st.GetUpVotes(),
			expectedResult: int64(10),
		},
		{
			name:           "test get down votes",
			actualResult:   st.GetDownVotes(),
			expectedResult: int64(2),
		},
		{
			name:           "test get created at",
			actualResult:   st.GetCreatedAt(),
			expectedResult: createdAt,
		},
		{
			name:           "test get updated at",
			actualResult:   st.GetUpdatedAt(),
			expectedResult: updatedAt,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedResult, testCase.actualResult)
		})
	}
}

func TestStoryAddView(t *testing.T) {
	id := "ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a"
	createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		actualResult   func() int64
		expectedResult int64
	}{
		{
			name: "test add views",
			actualResult: func() int64 {
				st, err := model.NewStoryBuilder().
					SetID(id).
					SetTitle(100, "title").
					SetBody(10000, "this is a test body").
					SetViewCount(25).
					SetUpVotes(10).
					SetDownVotes(2).
					SetCreatedAt(createdAt).
					SetUpdatedAt(updatedAt).
					Build()

				require.NoError(t, err)

				for i := 0; i < 10; i++ {
					st.AddView()
				}

				return st.GetViewCount()
			},
			expectedResult: int64(35),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedResult, testCase.actualResult())
		})
	}
}

func TestStoryUpVote(t *testing.T) {
	id := "ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a"
	createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		actualResult   func() int64
		expectedResult int64
	}{

		{
			name: "test add up votes",
			actualResult: func() int64 {
				st, err := model.NewStoryBuilder().
					SetID(id).
					SetTitle(100, "title").
					SetBody(10000, "this is a test body").
					SetViewCount(25).
					SetUpVotes(10).
					SetDownVotes(2).
					SetCreatedAt(createdAt).
					SetUpdatedAt(updatedAt).
					Build()

				require.NoError(t, err)

				for i := 0; i < 10; i++ {
					st.UpVote()
				}

				return st.GetUpVotes()
			},
			expectedResult: int64(20),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedResult, testCase.actualResult())
		})
	}
}

func TestStoryDownVote(t *testing.T) {
	id := "ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a"
	createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		actualResult   func() int64
		expectedResult int64
	}{
		{
			name: "test add down votes",
			actualResult: func() int64 {
				st, err := model.NewStoryBuilder().
					SetID(id).
					SetTitle(100, "title").
					SetBody(10000, "this is a test body").
					SetViewCount(25).
					SetUpVotes(10).
					SetDownVotes(2).
					SetCreatedAt(createdAt).
					SetUpdatedAt(updatedAt).
					Build()

				require.NoError(t, err)

				for i := 0; i < 10; i++ {
					st.DownVote()
				}

				return st.GetDownVotes()
			},
			expectedResult: int64(12),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedResult, testCase.actualResult())
		})
	}
}
