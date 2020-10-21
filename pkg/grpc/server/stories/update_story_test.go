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
	"strings"
	"testing"
	"time"
)

func TestStoriesServerUpdateStory(t *testing.T) {
	cfg := config.NewConfig("../../../../local.env").StoryConfig()

	createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

	testCases := map[string]struct {
		input          func() (service.StoryService, *proto.UpdateStoryRequest)
		expectedResult *proto.UpdateStoryResponse
		expectedError  error
	}{
		"test update story success": {
			input: func() (service.StoryService, *proto.UpdateStoryRequest) {
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

				return ms, req
			},
			expectedResult: &proto.UpdateStoryResponse{Success: true},
		},
		"test update story failure when uuid is invalid": {
			input: func() (service.StoryService, *proto.UpdateStoryRequest) {
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

				return &service.MockStoriesService{}, req
			},
			expectedResult: &proto.UpdateStoryResponse{Success: false},
			expectedError: liberr.WithArgs(
				liberr.Operation("Server.UpdateStory"),
				liberr.WithArgs(
					liberr.SeverityError,
					liberr.ValidationError,
					liberr.Operation("StoryBuilder.Build"),
					errors.New("invalid id: abc"),
				),
			),
		},
		"test update story failure when title is empty": {
			input: func() (service.StoryService, *proto.UpdateStoryRequest) {
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

				return &service.MockStoriesService{}, req
			},
			expectedResult: &proto.UpdateStoryResponse{Success: false},
			expectedError: liberr.WithArgs(
				liberr.Operation("Server.UpdateStory"),
				liberr.WithArgs(
					liberr.SeverityError,
					liberr.ValidationError,
					liberr.Operation("StoryBuilder.Build"),
					errors.New("title cannot be empty"),
				),
			),
		},
		"test update story failure when body is empty": {
			input: func() (service.StoryService, *proto.UpdateStoryRequest) {
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

				return &service.MockStoriesService{}, req
			},
			expectedResult: &proto.UpdateStoryResponse{Success: false},
			expectedError: liberr.WithArgs(
				liberr.Operation("Server.UpdateStory"),
				liberr.WithArgs(
					liberr.SeverityError,
					liberr.ValidationError,
					liberr.Operation("StoryBuilder.Build"),
					errors.New("body cannot be empty"),
				),
			),
		},
		"test update story failure when title exceeds max length": {
			input: func() (service.StoryService, *proto.UpdateStoryRequest) {
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

				return &service.MockStoriesService{}, req
			},
			expectedResult: &proto.UpdateStoryResponse{Success: false},
			expectedError: liberr.WithArgs(
				liberr.Operation("Server.UpdateStory"),
				liberr.WithArgs(
					liberr.SeverityError,
					liberr.ValidationError,
					liberr.Operation("StoryBuilder.Build"),
					errors.New("title max length exceeded"),
				),
			),
		},
		"test update story failure when body exceeds max length": {
			input: func() (service.StoryService, *proto.UpdateStoryRequest) {
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

				return &service.MockStoriesService{}, req
			},
			expectedResult: &proto.UpdateStoryResponse{Success: false},
			expectedError: liberr.WithArgs(
				liberr.Operation("Server.UpdateStory"),
				liberr.WithArgs(
					liberr.SeverityError,
					liberr.ValidationError,
					liberr.Operation("StoryBuilder.Build"),
					errors.New("body max length exceeded"),
				),
			),
		},
		"test update story failure when service returns error": {
			input: func() (service.StoryService, *proto.UpdateStoryRequest) {
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
				ms.On("UpdateStory", st).Return(int64(0), liberr.WithArgs(errors.New("failed to update story")))

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

				return ms, req
			},
			expectedResult: &proto.UpdateStoryResponse{Success: false},
			expectedError:  liberr.WithArgs(liberr.Operation("Server.UpdateStory"), liberr.WithArgs(errors.New("failed to update story"))),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			svc, req := testCase.input()

			testStoriesServerUpdateStory(t, testCase.expectedError, testCase.expectedResult, req, svc)
		})
	}
}

func testStoriesServerUpdateStory(t *testing.T, expectedError error, expectedResult *proto.UpdateStoryResponse, req *proto.UpdateStoryRequest, svc service.StoryService) {
	cfg := config.NewConfig("../../../../local.env").StoryConfig()

	server := stories.NewStoriesServer(cfg, svc)

	res, err := server.UpdateStory(context.Background(), req)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, expectedResult, res)
}
