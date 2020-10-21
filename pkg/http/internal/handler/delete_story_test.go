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
	"github.com/nsnikhil/stories/pkg/story/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteStory(t *testing.T) {
	testCases := map[string]struct {
		input          func() (service.StoryService, io.Reader)
		expectedResult string
		expectedCode   int
	}{
		"test delete story success": {
			input: func() (service.StoryService, io.Reader) {
				id := "adbca278-7e5c-4831-bf90-15fadfda0dd1"

				dlReq := contract.DeleteStoryRequest{StoryID: id}

				b, err := json.Marshal(dlReq)
				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("DeleteStory", id).Return(int64(1), nil)

				return ms, bytes.NewBuffer(b)
			},
			expectedCode:   http.StatusOK,
			expectedResult: "{\"data\":{\"success\":true},\"success\":true}",
		},
		"test delete story failure when req body is nil": {
			input: func() (service.StoryService, io.Reader) {
				return &service.MockStoriesService{}, nil
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"message\":\"unexpected end of JSON input\"},\"success\":false}",
		},
		"test delete story failure when svc call fails": {
			input: func() (service.StoryService, io.Reader) {
				id := "adbca278-7e5c-4831-bf90-15fadfda0dd1"

				dlReq := contract.DeleteStoryRequest{StoryID: id}

				b, err := json.Marshal(dlReq)
				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("DeleteStory", id).Return(int64(0), liberr.WithArgs(liberr.SeverityError, errors.New("failed to delete story")))

				return ms, bytes.NewBuffer(b)
			},
			expectedCode:   http.StatusInternalServerError,
			expectedResult: "{\"error\":{\"message\":\"internal server error\"},\"success\":false}",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			svc, body := testCase.input()

			testDeleteStory(t, testCase.expectedCode, testCase.expectedResult, svc, body)
		})
	}
}

func testDeleteStory(t *testing.T, expectedCode int, expectedBody string, svc service.StoryService, body io.Reader) {
	dh := handler.NewDeleteStoryHandler(svc)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/story/delete", body)

	mdl.WithError(reporters.NewLogger("dev", "debug"), dh.DeleteStory)(w, r)

	assert.Equal(t, expectedCode, w.Code)
	assert.Equal(t, expectedBody, w.Body.String())
}
