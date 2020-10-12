package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/nsnikhil/stories/pkg/http/internal/contract"
	"github.com/nsnikhil/stories/pkg/http/internal/handler"
	mdl "github.com/nsnikhil/stories/pkg/http/internal/middleware"
	"github.com/nsnikhil/stories/pkg/liberr"
	reporters "github.com/nsnikhil/stories/pkg/reporting"
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/nsnikhil/stories/pkg/story/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetStory(t *testing.T) {
	lgr := reporters.NewLogger("dev", "debug")

	testCases := []struct {
		name           string
		actualResult   func() (string, int)
		expectedResult string
		expectedCode   int
	}{
		{
			name: "test get story success",
			actualResult: func() (string, int) {
				id := "adbca278-7e5c-4831-bf90-15fadfda0dd1"

				createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

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
				ms.On("GetStory", id).Return(ds, nil)

				gtReq := contract.GetStoryRequest{
					StoryID: id,
				}

				b, err := json.Marshal(gtReq)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/story/get", bytes.NewBuffer(b))

				gh := handler.NewGetStoryHandler(ms)

				mdl.WithError(lgr, gh.GetStory)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusOK,
			expectedResult: "{\"data\":{\"id\":\"adbca278-7e5c-4831-bf90-15fadfda0dd1\",\"title\":\"title\",\"body\":\"test body\",\"view_count\":25,\"up_votes\":10,\"down_votes\":2,\"created_at\":1596038400,\"updated_at\":1596038400},\"success\":true}",
		},
		{
			name: "test get story failure when req body is nil",
			actualResult: func() (string, int) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/story/get", nil)

				gh := handler.NewGetStoryHandler(&service.MockStoriesService{})

				mdl.WithError(lgr, gh.GetStory)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"message\":\"unexpected end of JSON input\"},\"success\":false}",
		},
		{
			name: "test get story failure when service calls fails",
			actualResult: func() (string, int) {
				id := "adbca278-7e5c-4831-bf90-15fadfda0dd1"

				ms := &service.MockStoriesService{}
				ms.On("GetStory", id).Return(&model.Story{}, liberr.WithArgs(liberr.SeverityError, errors.New("failed to get story")))

				gtReq := contract.GetStoryRequest{
					StoryID: id,
				}

				b, err := json.Marshal(gtReq)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/story/get", bytes.NewBuffer(b))

				gh := handler.NewGetStoryHandler(ms)

				mdl.WithError(lgr, gh.GetStory)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusInternalServerError,
			expectedResult: "{\"error\":{\"message\":\"internal server error\"},\"success\":false}",
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
