package domain

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateNewStory(t *testing.T) {
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
			expectedError: errors.New("title is empty"),
		},
		{
			name: "test return error when body is empty for vanilla story",
			actualResult: func() (*Story, error) {
				return NewVanillaStory("title", "")
			},
			expectedError: errors.New("body is empty"),
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