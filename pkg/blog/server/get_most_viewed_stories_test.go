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

func TestStoriesServerGetMostViewedStories(t *testing.T) {
	cfg := config.LoadConfigs()
	lgr := zap.NewExample()
	createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		actualResult   func() (*proto.MostViewedStoriesResponse, error)
		expectedResult *proto.MostViewedStoriesResponse
		expectedError  error
	}{
		{
			name: "test get most viewed stories success",
			actualResult: func() (*proto.MostViewedStoriesResponse, error) {
				st, err := domain.NewStory("adbca278-7e5c-4831-bf90-15fadfda0dd1", "title", "test body", 25, 10, 2, createdAt, updatedAt)
				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("GetMostViewsStories", 0, 10).Return([]domain.Story{*st}, nil)

				deps := newDeps(service.NewService(ms), cfg, lgr)
				server := newStoriesServer(deps)

				req := &proto.MostViewedStoriesRequest{
					Offset: 0,
					Limit:  10,
				}

				return server.GetMostViewedStories(context.Background(), req)
			},
			expectedResult: &proto.MostViewedStoriesResponse{
				Stories: []*proto.Story{
					{
						Id:            "adbca278-7e5c-4831-bf90-15fadfda0dd1",
						Title:         "title",
						Body:          "test body",
						Views:         25,
						UpVotes:       10,
						DownVotes:     2,
						CreatedAtUnix: createdAt.Unix(),
						UpdatedAtUnix: updatedAt.Unix(),
					},
				},
			},
		},
		{
			name: "test get most viewed stories failure",
			actualResult: func() (*proto.MostViewedStoriesResponse, error) {
				ms := &service.MockStoriesService{}
				ms.On("GetMostViewsStories", 0, 10).Return([]domain.Story{}, errors.New("failed to get most viewed stories"))

				deps := newDeps(service.NewService(ms), cfg, lgr)
				server := newStoriesServer(deps)

				req := &proto.MostViewedStoriesRequest{
					Offset: 0,
					Limit:  10,
				}

				return server.GetMostViewedStories(context.Background(), req)
			},
			expectedResult: (*proto.MostViewedStoriesResponse)(nil),
			expectedError:  errors.New("failed to get most viewed stories"),
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
