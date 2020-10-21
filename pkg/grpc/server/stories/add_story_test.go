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
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestStoriesServerAddStory(t *testing.T) {
	testCases := map[string]struct {
		input          func() (service.StoryService, *proto.AddStoryRequest)
		expectedResult *proto.AddStoryResponse
		expectedError  error
	}{
		"test add story success": {
			input: func() (service.StoryService, *proto.AddStoryRequest) {
				ms := &service.MockStoriesService{}
				ms.On("AddStory", mock.AnythingOfType("*model.Story")).Return(nil)

				req := &proto.AddStoryRequest{
					Story: &proto.Story{
						Title: "title",
						Body:  "test body",
					},
				}

				return ms, req
			},
			expectedResult: &proto.AddStoryResponse{
				Success: true,
			},
		},
		"test add story return error when title is empty": {
			input: func() (service.StoryService, *proto.AddStoryRequest) {
				req := &proto.AddStoryRequest{
					Story: &proto.Story{
						Title: "",
						Body:  "test body",
					},
				}

				return &service.MockStoriesService{}, req
			},
			expectedResult: &proto.AddStoryResponse{
				Success: false,
			},
			expectedError: liberr.WithArgs(
				liberr.Operation("Server.AddStory"),
				liberr.WithArgs(
					liberr.SeverityError,
					liberr.ValidationError,
					liberr.Operation("StoryBuilder.Build"),
					errors.New("title cannot be empty"),
				),
			),
		},
		"test add story return error when body is empty": {
			input: func() (service.StoryService, *proto.AddStoryRequest) {
				req := &proto.AddStoryRequest{
					Story: &proto.Story{
						Title: "title",
						Body:  "",
					},
				}

				return &service.MockStoriesService{}, req
			},
			expectedResult: &proto.AddStoryResponse{
				Success: false,
			},
			expectedError: liberr.WithArgs(
				liberr.Operation("Server.AddStory"),
				liberr.WithArgs(
					liberr.SeverityError,
					liberr.ValidationError,
					liberr.Operation("StoryBuilder.Build"),
					errors.New("body cannot be empty"),
				),
			),
		},
		"test add story return error service calls fails": {
			input: func() (service.StoryService, *proto.AddStoryRequest) {
				ms := &service.MockStoriesService{}
				ms.On("AddStory", mock.AnythingOfType("*model.Story")).Return(liberr.WithArgs(errors.New("failed to add story")))

				req := &proto.AddStoryRequest{
					Story: &proto.Story{
						Title: "title",
						Body:  "test body",
					},
				}

				return ms, req
			},
			expectedResult: &proto.AddStoryResponse{
				Success: false,
			},
			expectedError: liberr.WithArgs(liberr.Operation("Server.AddStory"), liberr.WithArgs(errors.New("failed to add story"))),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			svc, req := testCase.input()
			testStoriesServerAddStory(t, testCase.expectedError, testCase.expectedResult, req, svc)
		})
	}
}

func testStoriesServerAddStory(t *testing.T, expectedError error, expectedResult *proto.AddStoryResponse, req *proto.AddStoryRequest, svc service.StoryService) {
	cfg := config.NewConfig("../../../../local.env").StoryConfig()

	server := stories.NewStoriesServer(cfg, svc)

	res, err := server.AddStory(context.Background(), req)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, expectedResult, res)
}
