package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/nsnikhil/stories/pkg/http/contract"
	"github.com/nsnikhil/stories/pkg/http/internal/handler"
	mdl "github.com/nsnikhil/stories/pkg/http/internal/middleware"
	"github.com/nsnikhil/stories/pkg/liberr"
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/nsnikhil/stories/pkg/story/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetTopRatedStories(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (string, int)
		expectedResult string
		expectedCode   int
	}{
		{
			name: "test get top rated stories success",
			actualResult: func() (string, int) {
				o, l := 10, 0

				tpReq := contract.TopRatedStoriesRequest{
					OffSet: o,
					Limit:  l,
				}

				b, err := json.Marshal(tpReq)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/story/top-rated", bytes.NewBuffer(b))

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
				ms.On("GetTopRatedStories", o, l).Return([]model.Story{*st}, nil)

				mvh := handler.NewGetTopRatedStoriesHandler(ms)

				mdl.WithError(mvh.GetTopRatedStories)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusOK,
			expectedResult: "{\"data\":[{\"id\":\"adbca278-7e5c-4831-bf90-15fadfda0dd1\",\"title\":\"title\",\"body\":\"test body\",\"view_count\":25,\"up_votes\":10,\"down_votes\":2,\"created_at\":1596038400,\"updated_at\":1596038400}],\"success\":true}",
		},
		{
			name: "test get top rated stories failure when req body is nil",
			actualResult: func() (string, int) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/story/top-rated", nil)

				mvh := handler.NewGetTopRatedStoriesHandler(&service.MockStoriesService{})

				mdl.WithError(mvh.GetTopRatedStories)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"message\":\"unexpected end of JSON input\"},\"success\":false}",
		},
		{
			name: "test get top rated stories failure when svc call fails",
			actualResult: func() (string, int) {
				o, l := 10, 0

				mvReq := contract.TopRatedStoriesRequest{
					OffSet: o,
					Limit:  l,
				}

				b, err := json.Marshal(mvReq)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/story/top-rated", bytes.NewBuffer(b))

				ms := &service.MockStoriesService{}
				ms.On("GetTopRatedStories", o, l).Return([]model.Story{}, liberr.WithArgs(liberr.SeverityError, errors.New("failed to get top rated stories")))

				mvh := handler.NewGetTopRatedStoriesHandler(ms)

				mdl.WithError(mvh.GetTopRatedStories)(w, r)

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
