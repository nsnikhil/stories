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
	"time"
)

func TestStoriesServerGetStory(t *testing.T) {
	cfg := config.LoadConfigs()
	lgr := zap.NewExample()
	createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		actualResult   func() (*proto.GetStoryResponse, error)
		expectedResult func() *proto.GetStoryResponse
		expectedError  error
	}{
		{
			name: "test get story success",
			actualResult: func() (*proto.GetStoryResponse, error) {
				id := "adbca278-7e5c-4831-bf90-15fadfda0dd1"

				st, err := domain.NewStory("adbca278-7e5c-4831-bf90-15fadfda0dd1", "title", "test body", 25, 10, 2, createdAt, updatedAt)
				require.NoError(t, err)
				st.Id = id

				ms := &service.MockStoriesService{}
				ms.On("GetStory", id).Return(st, nil)

				req := &proto.GetStoryRequest{
					StoryID: id,
				}

				deps := newDeps(service.NewService(ms), cfg, lgr)
				server := newStoriesServer(deps)

				return server.GetStory(context.Background(), req)
			},
			expectedResult: func() *proto.GetStoryResponse {
				return &proto.GetStoryResponse{
					Story: &proto.Story{
						Id:        "adbca278-7e5c-4831-bf90-15fadfda0dd1",
						Title:     "title",
						Body:      "test body",
						Views:     25,
						UpVotes:   10,
						DownVotes: 2,
						CreatedAtUnix: createdAt.Unix(),
						UpdatedAtUnix: updatedAt.Unix(),
					},
				}
			},
		},
		{
			name: "test get story failure",
			actualResult: func() (*proto.GetStoryResponse, error) {
				id := "adbca278-7e5c-4831-bf90-15fadfda0dd1"

				ms := &service.MockStoriesService{}
				ms.On("GetStory", id).Return(&domain.Story{}, errors.New("failed to get story"))

				req := &proto.GetStoryRequest{
					StoryID: id,
				}

				deps := newDeps(service.NewService(ms), cfg, lgr)
				server := newStoriesServer(deps)

				return server.GetStory(context.Background(), req)
			},
			expectedResult: func() *proto.GetStoryResponse {
				return (*proto.GetStoryResponse)(nil)
			},
			expectedError: errors.New("failed to get story"),
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