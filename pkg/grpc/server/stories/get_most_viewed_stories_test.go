package stories_test

import (
	"context"
	"errors"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/grpc/server/stories"
	"github.com/nsnikhil/stories/pkg/liberr"
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/nsnikhil/stories/pkg/story/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestStoriesServerGetMostViewedStories(t *testing.T) {
	cfg := config.NewConfig("../../../../local.env").StoryConfig()

	createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

	testCases := map[string]struct {
		input          func() service.StoryService
		expectedResult *proto.MostViewedStoriesResponse
		expectedError  error
	}{
		"test get most viewed story success": {
			input: func() service.StoryService {
				st, err := model.NewStoryBuilder().
					SetID("adbca278-7e5c-4831-bf90-15fadfda0dd1").
					SetTitle(cfg.TitleMaxLength(), "title").
					SetBody(cfg.BodyMaxLength(), "test body").
					SetViewCount(25).
					SetUpVotes(10).
					SetDownVotes(2).
					SetCreatedAt(createdAt).
					SetUpdatedAt(updatedAt).
					Build()

				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("GetMostViewsStories", 0, 10).Return([]model.Story{*st}, nil)

				return ms
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
		"test get most viewed story failure": {
			input: func() service.StoryService {
				ms := &service.MockStoriesService{}
				ms.On("GetMostViewsStories", 0, 10).Return([]model.Story{}, liberr.WithArgs(errors.New("failed to get most viewed story")))

				return ms
			},
			expectedResult: (*proto.MostViewedStoriesResponse)(nil),
			expectedError:  liberr.WithArgs(liberr.Operation("Server.GetMostViewedStories"), liberr.WithArgs(errors.New("failed to get most viewed story"))),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			svc := testCase.input()

			testStoriesServerGetMostViewedStories(t, testCase.expectedError, testCase.expectedResult, svc)
		})
	}
}

func testStoriesServerGetMostViewedStories(t *testing.T, expectedError error, expectedResult *proto.MostViewedStoriesResponse, svc service.StoryService) {
	cfg := config.NewConfig("../../../../local.env").StoryConfig()

	server := stories.NewStoriesServer(cfg, svc)

	req := &proto.MostViewedStoriesRequest{Offset: 0, Limit: 10}
	res, err := server.GetMostViewedStories(context.Background(), req)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, expectedResult, res)
}
