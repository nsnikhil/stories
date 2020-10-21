package model_test

import (
	"errors"
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/stretchr/testify/assert"
	"testing"
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
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, testCase.expectedResult, res)
		})
	}
}
