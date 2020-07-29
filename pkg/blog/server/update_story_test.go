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
	"strings"
	"testing"
	"time"
)

func TestStoriesServerUpdateStory(t *testing.T) {
	cfg := config.LoadConfigs()
	lgr := zap.NewExample()
	createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		actualResult   func() (*proto.UpdateStoryResponse, error)
		expectedResult *proto.UpdateStoryResponse
		expectedError  error
	}{
		{
			name: "test update story success",
			actualResult: func() (*proto.UpdateStoryResponse, error) {
				st, err := domain.NewStory("adbca278-7e5c-4831-bf90-15fadfda0dd1", "title", "test body", 25, 10, 2, createdAt, updatedAt)
				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("UpdateStory", st).Return(int64(1), nil)

				req := &proto.UpdateStoryRequest{
					Story: &proto.Story{
						Id:            "adbca278-7e5c-4831-bf90-15fadfda0dd1",
						Title:         "title",
						Body:          "test body",
						Views:         25,
						UpVotes:       10,
						DownVotes:     2,
						CreatedAtUnix: createdAt.Unix(),
						UpdatedAtUnix: updatedAt.Unix(),
					},
				}

				deps := newDeps(service.NewService(ms), cfg, lgr)
				server := newStoriesServer(deps)

				return server.UpdateStory(context.Background(), req)
			},
			expectedResult: &proto.UpdateStoryResponse{Success: true},
		},
		{
			name: "test update story failure when uuid is invalid",
			actualResult: func() (*proto.UpdateStoryResponse, error) {
				req := &proto.UpdateStoryRequest{
					Story: &proto.Story{
						Id:            "abc",
						Title:         "title",
						Body:          "test body",
						Views:         25,
						UpVotes:       10,
						DownVotes:     2,
						CreatedAtUnix: createdAt.Unix(),
						UpdatedAtUnix: updatedAt.Unix(),
					},
				}

				deps := newDeps(service.NewService(&service.MockStoriesService{}), cfg, lgr)
				server := newStoriesServer(deps)

				return server.UpdateStory(context.Background(), req)
			},
			expectedResult: &proto.UpdateStoryResponse{Success: false},
			expectedError:  errors.New("invalid id: abc"),
		},
		{
			name: "test update story failure when title is empty",
			actualResult: func() (*proto.UpdateStoryResponse, error) {
				req := &proto.UpdateStoryRequest{
					Story: &proto.Story{
						Id:            "adbca278-7e5c-4831-bf90-15fadfda0dd1",
						Title:         "",
						Body:          "test body",
						Views:         25,
						UpVotes:       10,
						DownVotes:     2,
						CreatedAtUnix: createdAt.Unix(),
						UpdatedAtUnix: updatedAt.Unix(),
					},
				}

				deps := newDeps(service.NewService(&service.MockStoriesService{}), cfg, lgr)
				server := newStoriesServer(deps)

				return server.UpdateStory(context.Background(), req)
			},
			expectedResult: &proto.UpdateStoryResponse{Success: false},
			expectedError:  errors.New("title cannot be empty"),
		},
		{
			name: "test update story failure when body is empty",
			actualResult: func() (*proto.UpdateStoryResponse, error) {
				req := &proto.UpdateStoryRequest{
					Story: &proto.Story{
						Id:            "adbca278-7e5c-4831-bf90-15fadfda0dd1",
						Title:         "title",
						Body:          "",
						Views:         25,
						UpVotes:       10,
						DownVotes:     2,
						CreatedAtUnix: createdAt.Unix(),
						UpdatedAtUnix: updatedAt.Unix(),
					},
				}

				deps := newDeps(service.NewService(&service.MockStoriesService{}), cfg, lgr)
				server := newStoriesServer(deps)

				return server.UpdateStory(context.Background(), req)
			},
			expectedResult: &proto.UpdateStoryResponse{Success: false},
			expectedError:  errors.New("body cannot be empty"),
		},
		{
			name: "test update story failure when title exceeds max length",
			actualResult: func() (*proto.UpdateStoryResponse, error) {
				var title strings.Builder
				for i := 0; i < 101; i++ {
					title.WriteString("a")
				}

				req := &proto.UpdateStoryRequest{
					Story: &proto.Story{
						Id:            "adbca278-7e5c-4831-bf90-15fadfda0dd1",
						Title:         title.String(),
						Body:          "test body",
						Views:         25,
						UpVotes:       10,
						DownVotes:     2,
						CreatedAtUnix: createdAt.Unix(),
						UpdatedAtUnix: updatedAt.Unix(),
					},
				}

				deps := newDeps(service.NewService(&service.MockStoriesService{}), cfg, lgr)
				server := newStoriesServer(deps)

				return server.UpdateStory(context.Background(), req)
			},
			expectedResult: &proto.UpdateStoryResponse{Success: false},
			expectedError:  errors.New("title max length exceeded"),
		},
		{
			name: "test update story failure when body exceeds max length",
			actualResult: func() (*proto.UpdateStoryResponse, error) {
				var body strings.Builder
				for i := 0; i < 100001; i++ {
					body.WriteString("a")
				}

				req := &proto.UpdateStoryRequest{
					Story: &proto.Story{
						Id:            "adbca278-7e5c-4831-bf90-15fadfda0dd1",
						Title:         "title",
						Body:          body.String(),
						Views:         25,
						UpVotes:       10,
						DownVotes:     2,
						CreatedAtUnix: createdAt.Unix(),
						UpdatedAtUnix: updatedAt.Unix(),
					},
				}

				deps := newDeps(service.NewService(&service.MockStoriesService{}), cfg, lgr)
				server := newStoriesServer(deps)

				return server.UpdateStory(context.Background(), req)
			},
			expectedResult: &proto.UpdateStoryResponse{Success: false},
			expectedError:  errors.New("body max length exceeded"),
		},
		{
			name: "test update story failure when service returns error",
			actualResult: func() (*proto.UpdateStoryResponse, error) {
				st, err := domain.NewStory("adbca278-7e5c-4831-bf90-15fadfda0dd1", "title", "test body", 25, 10, 2, createdAt, updatedAt)
				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("UpdateStory", st).Return(int64(0), errors.New("failed to update story"))

				req := &proto.UpdateStoryRequest{
					Story: &proto.Story{
						Id:            "adbca278-7e5c-4831-bf90-15fadfda0dd1",
						Title:         "title",
						Body:          "test body",
						Views:         25,
						UpVotes:       10,
						DownVotes:     2,
						CreatedAtUnix: createdAt.Unix(),
						UpdatedAtUnix: updatedAt.Unix(),
					},
				}

				deps := newDeps(service.NewService(ms), cfg, lgr)
				server := newStoriesServer(deps)

				return server.UpdateStory(context.Background(), req)
			},
			expectedResult: &proto.UpdateStoryResponse{Success: false},
			expectedError:  errors.New("failed to update story"),
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
