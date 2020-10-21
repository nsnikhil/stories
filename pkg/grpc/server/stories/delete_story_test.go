package stories_test

import (
	"context"
	"errors"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/grpc/server/stories"
	"github.com/nsnikhil/stories/pkg/liberr"
	"github.com/nsnikhil/stories/pkg/story/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStoriesServerDeleteStory(t *testing.T) {

	testCases := map[string]struct {
		input          func() service.StoryService
		expectedResult *proto.DeleteStoryResponse
		expectedError  error
	}{
		"test delete story success": {
			input: func() service.StoryService {
				ms := &service.MockStoriesService{}
				ms.On("DeleteStory", "adbca278-7e5c-4831-bf90-15fadfda0dd1").Return(int64(1), nil)
				return ms
			},
			expectedResult: &proto.DeleteStoryResponse{
				Success: true,
			},
		},
		"test delete story service failure": {
			input: func() service.StoryService {
				ms := &service.MockStoriesService{}
				ms.On("DeleteStory", "adbca278-7e5c-4831-bf90-15fadfda0dd1").Return(int64(0), liberr.WithArgs(errors.New("failed to delete story")))
				return ms
			},
			expectedResult: &proto.DeleteStoryResponse{
				Success: false,
			},
			expectedError: liberr.WithArgs(liberr.Operation("Server.DeleteStory"), liberr.WithArgs(errors.New("failed to delete story"))),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			svc := testCase.input()

			testStoriesServerDeleteStory(t, testCase.expectedError, testCase.expectedResult, svc)
		})
	}
}

func testStoriesServerDeleteStory(t *testing.T, expectedError error, expectedResult *proto.DeleteStoryResponse, svc service.StoryService) {
	cfg := config.NewConfig("../../../../local.env").StoryConfig()
	id := "adbca278-7e5c-4831-bf90-15fadfda0dd1"

	req := &proto.DeleteStoryRequest{StoryID: id}

	server := stories.NewStoriesServer(cfg, svc)

	res, err := server.DeleteStory(context.Background(), req)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, expectedResult, res)
}
