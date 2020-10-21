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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetMostViewedStories(t *testing.T) {
	testCases := map[string]struct {
		input          func() (service.StoryService, io.Reader)
		expectedResult string
		expectedCode   int
	}{
		"test get most viewed stories success": {
			input: func() (service.StoryService, io.Reader) {
				o, l := 10, 0

				mvReq := contract.MostViewedStoriesRequest{OffSet: o, Limit: l}
				b, err := json.Marshal(mvReq)
				require.NoError(t, err)

				createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

				st, err := model.NewStoryBuilder().
					SetID("adbca278-7e5c-4831-bf90-15fadfda0dd1").
					SetTitle(10, "title").
					SetBody(10, "test body").
					SetViewCount(25).
					SetUpVotes(10).
					SetDownVotes(2).
					SetCreatedAt(createdAt).
					SetUpdatedAt(updatedAt).
					Build()

				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("GetMostViewsStories", o, l).Return([]model.Story{*st}, nil)

				return ms, bytes.NewBuffer(b)
			},
			expectedCode:   http.StatusOK,
			expectedResult: "{\"data\":[{\"id\":\"adbca278-7e5c-4831-bf90-15fadfda0dd1\",\"title\":\"title\",\"body\":\"test body\",\"view_count\":25,\"up_votes\":10,\"down_votes\":2,\"created_at\":1596038400,\"updated_at\":1596038400}],\"success\":true}",
		},
		"test get most viewed stories failure when req body is nil": {
			input: func() (service.StoryService, io.Reader) {
				return &service.MockStoriesService{}, nil
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"message\":\"unexpected end of JSON input\"},\"success\":false}",
		},
		"test get most viewed stories failure when svc call fails": {
			input: func() (service.StoryService, io.Reader) {
				o, l := 10, 0

				mvReq := contract.MostViewedStoriesRequest{OffSet: o, Limit: l}
				b, err := json.Marshal(mvReq)
				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("GetMostViewsStories", o, l).Return([]model.Story{}, liberr.WithArgs(liberr.SeverityError, errors.New("failed to get most viewed stories")))

				return ms, bytes.NewBuffer(b)
			},
			expectedCode:   http.StatusInternalServerError,
			expectedResult: "{\"error\":{\"message\":\"internal server error\"},\"success\":false}",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			svc, body := testCase.input()
			testGetMostViewedStories(t, testCase.expectedCode, testCase.expectedResult, svc, body)
		})
	}
}

func testGetMostViewedStories(t *testing.T, expectedCode int, expectedBody string, svc service.StoryService, body io.Reader) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/story/most-viewed", body)

	mvh := handler.NewGetMostViewedStoriesHandler(svc)

	mdl.WithError(reporters.NewLogger("dev", "debug"), mvh.GetMostViewedStories)(w, r)

	assert.Equal(t, expectedCode, w.Code)
	assert.Equal(t, expectedBody, w.Body.String())
}
