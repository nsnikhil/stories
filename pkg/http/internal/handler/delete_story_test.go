package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/nsnikhil/stories/pkg/http/contract"
	"github.com/nsnikhil/stories/pkg/http/internal/handler"
	mdl "github.com/nsnikhil/stories/pkg/http/internal/middleware"
	"github.com/nsnikhil/stories/pkg/story/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteStory(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (string, int)
		expectedResult string
		expectedCode   int
	}{
		{
			name: "test delete story success",
			actualResult: func() (string, int) {
				id := "adbca278-7e5c-4831-bf90-15fadfda0dd1"

				dlReq := contract.DeleteStoryRequest{
					StoryID: id,
				}

				b, err := json.Marshal(dlReq)
				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("DeleteStory", id).Return(int64(1), nil)

				dh := handler.NewDeleteStoryHandler(ms)

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodDelete, "/story/delete", bytes.NewBuffer(b))

				mdl.WithError(dh.DeleteStory)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusOK,
			expectedResult: "{\"data\":{\"success\":true},\"success\":true}",
		},
		{
			name: "test delete story failure when req body is nil",
			actualResult: func() (string, int) {
				dh := handler.NewDeleteStoryHandler(&service.MockStoriesService{})

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodDelete, "/story/delete", nil)

				mdl.WithError(dh.DeleteStory)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"code\":\"STx0001\",\"message\":\"unexpected end of JSON input\"},\"success\":false}",
		},
		{
			name: "test delete story failure when svc call fails",
			actualResult: func() (string, int) {
				id := "adbca278-7e5c-4831-bf90-15fadfda0dd1"

				dlReq := contract.DeleteStoryRequest{
					StoryID: id,
				}

				b, err := json.Marshal(dlReq)
				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("DeleteStory", id).Return(int64(0), errors.New("failed to delete story"))

				dh := handler.NewDeleteStoryHandler(ms)

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodDelete, "/story/delete", bytes.NewBuffer(b))

				mdl.WithError(dh.DeleteStory)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusInternalServerError,
			expectedResult: "{\"error\":{\"code\":\"STx0010\",\"message\":\"failed to delete story\"},\"success\":false}",
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
