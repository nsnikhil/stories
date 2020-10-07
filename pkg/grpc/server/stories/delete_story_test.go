package stories_test

import (
	"context"
	"errors"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/grpc/server/stories"
	"github.com/nsnikhil/stories/pkg/story/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStoriesServerDeleteStory(t *testing.T) {
	cfg := config.NewConfig("../../../../local.env").StoryConfig()

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

				ms := &service.MockStoriesService{}
				ms.On("DeleteStory", id).Return(int64(1), nil)

				req := &proto.DeleteStoryRequest{
					StoryID: id,
				}

				server := stories.NewStoriesServer(cfg, ms)
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

				ms := &service.MockStoriesService{}
				ms.On("DeleteStory", id).Return(int64(0), errors.New("failed to delete story"))

				req := &proto.DeleteStoryRequest{
					StoryID: id,
				}

				server := stories.NewStoriesServer(cfg, ms)
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
