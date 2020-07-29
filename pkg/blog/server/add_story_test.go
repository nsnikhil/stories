package server

import (
	"context"
	"errors"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/cmd/config"
	"github.com/nsnikhil/stories/pkg/blog/domain"
	"github.com/nsnikhil/stories/pkg/blog/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestStoriesServerAddStory(t *testing.T) {
	cfg := config.LoadConfigs()
	lgr := zap.NewExample()

	testCases := []struct {
		name           string
		actualResult   func() (*proto.AddStoryResponse, error)
		expectedResult *proto.AddStoryResponse
		expectedError  error
	}{
		{
			name: "test add story success",
			actualResult: func() (*proto.AddStoryResponse, error) {
				st, err := domain.NewVanillaStory("title", "test body")
				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("AddStory", st).Return(nil)

				req := &proto.AddStoryRequest{
					Story: &proto.Story{
						Title: "title",
						Body:  "test body",
					},
				}

				deps := newDeps(service.NewService(ms), cfg, lgr)
				server := newStoriesServer(deps)
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

				deps := newDeps(service.NewService(&service.MockStoriesService{}), cfg, lgr)
				server := newStoriesServer(deps)
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

				deps := newDeps(service.NewService(&service.MockStoriesService{}), cfg, lgr)
				server := newStoriesServer(deps)
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
				st, err := domain.NewVanillaStory("title", "test body")
				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("AddStory", st).Return(errors.New("failed to add story"))

				req := &proto.AddStoryRequest{
					Story: &proto.Story{
						Title: "title",
						Body:  "test body",
					},
				}

				deps := newDeps(service.NewService(ms), cfg, lgr)
				server := newStoriesServer(deps)
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
