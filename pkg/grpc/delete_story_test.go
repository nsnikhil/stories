package grpc

import (
	"context"
	"errors"
	"github.com/nsnikhil/stories-proto/proto"
	config2 "github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/story"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestStoriesServerDeleteStory(t *testing.T) {
	cfg := config2.LoadConfigs()
	lgr := zap.NewExample()

	testCases := []struct {
		name           string
		actualResult   func() (*proto.DeleteStoryResponse, error)
		expectedResult func() *proto.DeleteStoryResponse
		expectedError  error
	}{
		{
			name: "test delete story success",
			actualResult: func() (*proto.DeleteStoryResponse, error) {
				id := "adbca278-7e5c-4831-bf90-15fadfda0dd1"

				ms := &story.MockStoriesService{}
				ms.On("DeleteStory", id).Return(int64(1), nil)

				req := &proto.DeleteStoryRequest{
					StoryID: id,
				}

				deps := newDeps(story.NewService(ms), cfg, lgr)
				server := newStoriesServer(deps)

				return server.DeleteStory(context.Background(), req)
			},
			expectedResult: func() *proto.DeleteStoryResponse {
				return &proto.DeleteStoryResponse{
					Success: true,
				}
			},
		},
		{
			name: "test delete story service failure",
			actualResult: func() (*proto.DeleteStoryResponse, error) {
				id := "adbca278-7e5c-4831-bf90-15fadfda0dd1"

				ms := &story.MockStoriesService{}
				ms.On("DeleteStory", id).Return(int64(0), errors.New("failed to delete story"))

				req := &proto.DeleteStoryRequest{
					StoryID: id,
				}

				deps := newDeps(story.NewService(ms), cfg, lgr)
				server := newStoriesServer(deps)

				return server.DeleteStory(context.Background(), req)
			},
			expectedResult: func() *proto.DeleteStoryResponse {
				return &proto.DeleteStoryResponse{
					Success: false,
				}
			},
			expectedError: errors.New("failed to delete story"),
		},
		{
			name: "test delete story not deleted",
			actualResult: func() (*proto.DeleteStoryResponse, error) {
				id := "adbca278-7e5c-4831-bf90-15fadfda0dd1"

				ms := &story.MockStoriesService{}
				ms.On("DeleteStory", id).Return(int64(0), nil)

				req := &proto.DeleteStoryRequest{
					StoryID: id,
				}

				deps := newDeps(story.NewService(ms), cfg, lgr)
				server := newStoriesServer(deps)

				return server.DeleteStory(context.Background(), req)
			},
			expectedResult: func() *proto.DeleteStoryResponse {
				return &proto.DeleteStoryResponse{
					Success: false,
				}
			},
			expectedError: errors.New("failed to delete story"),
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
