package model_test

import (
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

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
