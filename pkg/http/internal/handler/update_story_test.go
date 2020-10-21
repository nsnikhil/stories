package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/http/internal/contract"
	"github.com/nsnikhil/stories/pkg/http/internal/handler"
	mdl "github.com/nsnikhil/stories/pkg/http/internal/middleware"
	"github.com/nsnikhil/stories/pkg/liberr"
	reporters "github.com/nsnikhil/stories/pkg/reporting"
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/nsnikhil/stories/pkg/story/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUpdateStory(t *testing.T) {
	testCases := map[string]struct {
		input          func() (service.StoryService, io.Reader)
		expectedResult string
		expectedCode   int
	}{
		"test update story success": {
			input: func() (service.StoryService, io.Reader) {
				createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

				upReq := contract.UpdateStoryRequest{
					Story: contract.Story{
						ID:        "adbca278-7e5c-4831-bf90-15fadfda0dd1",
						Title:     "title",
						Body:      "test body",
						ViewCount: 25,
						UpVotes:   10,
						DownVotes: 2,
						CreatedAt: createdAt.Unix(),
						UpdatedAt: updatedAt.Unix(),
					},
				}

				b, err := json.Marshal(upReq)
				require.NoError(t, err)

				ds, err := model.NewStoryBuilder().
					SetID("adbca278-7e5c-4831-bf90-15fadfda0dd1").
					SetTitle(100, "title").
					SetBody(100, "test body").
					SetViewCount(25).
					SetUpVotes(10).
					SetDownVotes(2).
					SetCreatedAt(createdAt).
					SetUpdatedAt(updatedAt).
					Build()

				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("UpdateStory", ds).Return(int64(1), nil)

				return ms, bytes.NewBuffer(b)
			},
			expectedCode:   http.StatusOK,
			expectedResult: "{\"data\":{\"success\":true},\"success\":true}",
		},
		"test update story failure when body is nil": {
			input: func() (service.StoryService, io.Reader) {
				return &service.MockStoriesService{}, nil
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"message\":\"unexpected end of JSON input\"},\"success\":false}",
		},
		"test update story failure when id is invalid": {
			input: func() (service.StoryService, io.Reader) {
				createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

				upReq := contract.UpdateStoryRequest{
					Story: contract.Story{
						ID:        "invalid-id",
						Title:     "title",
						Body:      "test body",
						ViewCount: 25,
						UpVotes:   10,
						DownVotes: 2,
						CreatedAt: createdAt.Unix(),
						UpdatedAt: updatedAt.Unix(),
					},
				}

				b, err := json.Marshal(upReq)
				require.NoError(t, err)

				return &service.MockStoriesService{}, bytes.NewBuffer(b)
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"message\":\"invalid id: invalid-id\"},\"success\":false}",
		},
		"test update story failure when title is empty": {
			input: func() (service.StoryService, io.Reader) {
				createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

				upReq := contract.UpdateStoryRequest{
					Story: contract.Story{
						ID:        "adbca278-7e5c-4831-bf90-15fadfda0dd1",
						Title:     "",
						Body:      "test body",
						ViewCount: 25,
						UpVotes:   10,
						DownVotes: 2,
						CreatedAt: createdAt.Unix(),
						UpdatedAt: updatedAt.Unix(),
					},
				}

				b, err := json.Marshal(upReq)
				require.NoError(t, err)

				return &service.MockStoriesService{}, bytes.NewBuffer(b)
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"message\":\"title cannot be empty\"},\"success\":false}",
		},
		"test update story failure when body is empty": {
			input: func() (service.StoryService, io.Reader) {
				createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

				upReq := contract.UpdateStoryRequest{
					Story: contract.Story{
						ID:        "adbca278-7e5c-4831-bf90-15fadfda0dd1",
						Title:     "title",
						Body:      "",
						ViewCount: 25,
						UpVotes:   10,
						DownVotes: 2,
						CreatedAt: createdAt.Unix(),
						UpdatedAt: updatedAt.Unix(),
					},
				}

				b, err := json.Marshal(upReq)
				require.NoError(t, err)

				return &service.MockStoriesService{}, bytes.NewBuffer(b)
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"message\":\"body cannot be empty\"},\"success\":false}",
		},
		"test update story failure when svc call fails": {
			input: func() (service.StoryService, io.Reader) {
				createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

				upReq := contract.UpdateStoryRequest{
					Story: contract.Story{
						ID:        "adbca278-7e5c-4831-bf90-15fadfda0dd1",
						Title:     "title",
						Body:      "test body",
						ViewCount: 25,
						UpVotes:   10,
						DownVotes: 2,
						CreatedAt: createdAt.Unix(),
						UpdatedAt: updatedAt.Unix(),
					},
				}

				b, err := json.Marshal(upReq)
				require.NoError(t, err)

				ds, err := model.NewStoryBuilder().
					SetID("adbca278-7e5c-4831-bf90-15fadfda0dd1").
					SetTitle(100, "title").
					SetBody(100, "test body").
					SetViewCount(25).
					SetUpVotes(10).
					SetDownVotes(2).
					SetCreatedAt(createdAt).
					SetUpdatedAt(updatedAt).
					Build()

				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("UpdateStory", ds).Return(int64(0), liberr.WithArgs(liberr.SeverityError, errors.New("failed to update story")))

				return ms, bytes.NewBuffer(b)
			},
			expectedCode:   http.StatusInternalServerError,
			expectedResult: "{\"error\":{\"message\":\"internal server error\"},\"success\":false}",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			svc, body := testCase.input()

			testUpdateStory(t, testCase.expectedCode, testCase.expectedResult, svc, body)
		})
	}
}

func testUpdateStory(t *testing.T, expectedCode int, expectedBody string, svc service.StoryService, body io.Reader) {
	cfg := config.NewConfig("../../../../local.env")

	uh := handler.NewUpdateStoryHandler(cfg.StoryConfig(), svc)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/story/update", body)

	mdl.WithError(reporters.NewLogger("dev", "debug"), uh.UpdateStory)(w, r)

	assert.Equal(t, expectedCode, w.Code)
	assert.Equal(t, expectedBody, w.Body.String())
}
