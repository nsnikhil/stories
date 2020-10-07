package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/http/contract"
	"github.com/nsnikhil/stories/pkg/http/internal/handler"
	mdl "github.com/nsnikhil/stories/pkg/http/internal/middleware"
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/nsnikhil/stories/pkg/story/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUpdateStory(t *testing.T) {
	cfg := config.NewConfig("../../../../local.env")

	testCases := []struct {
		name           string
		actualResult   func() (string, int)
		expectedResult string
		expectedCode   int
	}{
		{
			name: "test update story success",
			actualResult: func() (string, int) {
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

				uh := handler.NewUpdateStoryHandler(cfg.StoryConfig(), ms)

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPatch, "/story/update", bytes.NewBuffer(b))

				mdl.WithError(uh.UpdateStory)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusOK,
			expectedResult: "{\"data\":{\"success\":true},\"success\":true}",
		},
		{
			name: "test update story failure when body is nil",
			actualResult: func() (string, int) {
				uh := handler.NewUpdateStoryHandler(cfg.StoryConfig(), &service.MockStoriesService{})

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPatch, "/story/update", nil)

				mdl.WithError(uh.UpdateStory)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"code\":\"STx0001\",\"message\":\"unexpected end of JSON input\"},\"success\":false}",
		},
		{
			name: "test update story failure when id is invalid",
			actualResult: func() (string, int) {
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

				uh := handler.NewUpdateStoryHandler(cfg.StoryConfig(), &service.MockStoriesService{})

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPatch, "/story/update", bytes.NewBuffer(b))

				mdl.WithError(uh.UpdateStory)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"code\":\"STx0001\",\"message\":\"invalid id: invalid-id\"},\"success\":false}",
		},
		{
			name: "test update story failure when title is empty",
			actualResult: func() (string, int) {
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

				uh := handler.NewUpdateStoryHandler(cfg.StoryConfig(), &service.MockStoriesService{})

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPatch, "/story/update", bytes.NewBuffer(b))

				mdl.WithError(uh.UpdateStory)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"code\":\"STx0001\",\"message\":\"title cannot be empty\"},\"success\":false}",
		},
		{
			name: "test update story failure when body is empty",
			actualResult: func() (string, int) {
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

				uh := handler.NewUpdateStoryHandler(cfg.StoryConfig(), &service.MockStoriesService{})

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPatch, "/story/update", bytes.NewBuffer(b))

				mdl.WithError(uh.UpdateStory)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"code\":\"STx0001\",\"message\":\"body cannot be empty\"},\"success\":false}",
		},
		{
			name: "test update story failure when svc call fails",
			actualResult: func() (string, int) {
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
				ms.On("UpdateStory", ds).Return(int64(0), errors.New("failed to update story"))

				uh := handler.NewUpdateStoryHandler(cfg.StoryConfig(), ms)

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPatch, "/story/update", bytes.NewBuffer(b))

				mdl.WithError(uh.UpdateStory)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusInternalServerError,
			expectedResult: "{\"error\":{\"code\":\"STx0010\",\"message\":\"failed to update story\"},\"success\":false}",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, code := testCase.actualResult()

			assert.Equal(t, testCase.expectedCode, code)
			assert.Equal(t, testCase.expectedResult, res)
		})
	}
}
