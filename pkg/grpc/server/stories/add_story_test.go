package stories_test

import (
	"context"
	"errors"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/grpc/server/stories"
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/nsnikhil/stories/pkg/story/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStoriesServerAddStory(t *testing.T) {
	cfg := config.NewConfig("../../../../local.env").StoryConfig()

	testCases := []struct {
		name           string
		actualResult   func() (*proto.AddStoryResponse, error)
		expectedResult *proto.AddStoryResponse
		expectedError  error
	}{
		{
			name: "test add story success",
			actualResult: func() (*proto.AddStoryResponse, error) {
				st, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("AddStory", st).Return(nil)

				req := &proto.AddStoryRequest{
					Story: &proto.Story{
						Title: "title",
						Body:  "test body",
					},
				}

				server := stories.NewStoriesServer(cfg, ms)
				return server.AddStory(context.Background(), req)
			},
			expectedResult: &proto.AddStoryResponse{
				Success: true,
			},
		},
		{
			name: "test add story return error when title is empty",
			actualResult: func() (*proto.AddStoryResponse, error) {
				req := &proto.AddStoryRequest{
					Story: &proto.Story{
						Title: "",
						Body:  "test body",
					},
				}

				server := stories.NewStoriesServer(cfg, &service.MockStoriesService{})
				return server.AddStory(context.Background(), req)
			},
			expectedResult: &proto.AddStoryResponse{
				Success: false,
			},
			expectedError: errors.New("title cannot be empty"),
		},
		{
			name: "test add story return error when body is empty",
			actualResult: func() (*proto.AddStoryResponse, error) {
				req := &proto.AddStoryRequest{
					Story: &proto.Story{
						Title: "title",
						Body:  "",
					},
				}

				server := stories.NewStoriesServer(cfg, &service.MockStoriesService{})
				return server.AddStory(context.Background(), req)
			},
			expectedResult: &proto.AddStoryResponse{
				Success: false,
			},
			expectedError: errors.New("body cannot be empty"),
		},
		{
			name: "test add story return error service calls fails",
			actualResult: func() (*proto.AddStoryResponse, error) {
				st, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("AddStory", st).Return(errors.New("failed to add story"))

				req := &proto.AddStoryRequest{
					Story: &proto.Story{
						Title: "title",
						Body:  "test body",
					},
				}

				server := stories.NewStoriesServer(cfg, ms)
				return server.AddStory(context.Background(), req)
			},
			expectedResult: &proto.AddStoryResponse{
				Success: false,
			},
			expectedError: errors.New("failed to add story"),
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
