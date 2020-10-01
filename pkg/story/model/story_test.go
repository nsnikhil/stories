package model

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateNewVanillaStory(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (*Story, error)
		expectedResult *Story
		expectedError  error
	}{
		{
			name: "test create new vanilla story",
			actualResult: func() (*Story, error) {
				return NewVanillaStory("title", "this is a test body")
			},
			expectedResult: &Story{
				Title: "title",
				Body:  "this is a test body",
			},
		},
		{
			name: "test return error when title is empty for vanilla story",
			actualResult: func() (*Story, error) {
				return NewVanillaStory("", "this is a test body")
			},
			expectedError: errors.New("title cannot be empty"),
		},
		{
			name: "test return error when body is empty for vanilla story",
			actualResult: func() (*Story, error) {
				return NewVanillaStory("title", "")
			},
			expectedError: errors.New("body cannot be empty"),
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

func TestCreateNewStory(t *testing.T) {
	id := "ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a"
	createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		actualResult   func() (*Story, error)
		expectedResult *Story
		expectedError  error
	}{
		{
			name: "test create new story",
			actualResult: func() (*Story, error) {
				return NewStory(id, "title", "this is a test body", 0, 0, 0, createdAt, updatedAt)
			},
			expectedResult: &Story{
				ID:        id,
				Title:     "title",
				Body:      "this is a test body",
				ViewCount: 0,
				UpVotes:   0,
				DownVotes: 0,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
		},
		{
			name: "test return error when id is invalid",
			actualResult: func() (*Story, error) {
				return NewStory("invalid-id", "title", "this is a test body", 0, 0, 0, createdAt, updatedAt)
			},
			expectedError: errors.New("invalid id: invalid-id"),
		},
		{
			name: "test return error when title is empty for story",
			actualResult: func() (*Story, error) {
				return NewStory(id, "", "this is a test body", 0, 0, 0, createdAt, updatedAt)
			},
			expectedError: errors.New("title cannot be empty"),
		},
		{
			name: "test return error when body is empty for story",
			actualResult: func() (*Story, error) {
				return NewStory(id, "title", "", 0, 0, 0, createdAt, updatedAt)
			},
			expectedError: errors.New("body cannot be empty"),
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

func TestStoryGetter(t *testing.T) {
	id := "ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a"
	createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

	st, err := NewStory(id, "title", "this is a test body", 25, 10, 2, createdAt, updatedAt)
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
		{
			name:           "test get table name",
			actualResult:   st.TableName(),
			expectedResult: "story",
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
			name: "test add one view",
			actualResult: func() int64 {
				st, err := NewStory(id, "title", "this is a test body", 25, 10, 2, createdAt, updatedAt)
				require.NoError(t, err)

				st.AddView()

				return st.GetViewCount()
			},
			expectedResult: int64(26),
		},
		{
			name: "test add ten views",
			actualResult: func() int64 {
				st, err := NewStory(id, "title", "this is a test body", 25, 10, 2, createdAt, updatedAt)
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
			name: "test add one up vote",
			actualResult: func() int64 {
				st, err := NewStory(id, "title", "this is a test body", 25, 10, 2, createdAt, updatedAt)
				require.NoError(t, err)

				st.UpVote()

				return st.GetUpVotes()
			},
			expectedResult: int64(11),
		},
		{
			name: "test add ten up votes",
			actualResult: func() int64 {
				st, err := NewStory(id, "title", "this is a test body", 25, 10, 2, createdAt, updatedAt)
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
			name: "test add down vote",
			actualResult: func() int64 {
				st, err := NewStory(id, "title", "this is a test body", 25, 10, 2, createdAt, updatedAt)
				require.NoError(t, err)

				st.DownVote()

				return st.GetDownVotes()
			},
			expectedResult: int64(3),
		},
		{
			name: "test add ten down votes",
			actualResult: func() int64 {
				st, err := NewStory(id, "title", "this is a test body", 25, 10, 2, createdAt, updatedAt)
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
