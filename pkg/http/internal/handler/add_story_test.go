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
	"github.com/nsnikhil/stories/pkg/story/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddStory(t *testing.T) {
	testCases := map[string]struct {
		input          func() (service.StoryService, io.Reader)
		expectedResult string
		expectedCode   int
	}{
		"test add story success": {
			input: func() (service.StoryService, io.Reader) {
				ms := &service.MockStoriesService{}
				ms.On("AddStory", mock.AnythingOfType("*model.Story")).Return(nil)

				reqSt := contract.AddStoryRequest{Title: "title", Body: "test body"}
				b, err := json.Marshal(&reqSt)
				require.NoError(t, err)

				return ms, bytes.NewBuffer(b)
			},
			expectedCode:   http.StatusCreated,
			expectedResult: "{\"data\":{\"success\":true},\"success\":true}",
		},
		"test add story fails when body is nil": {
			input: func() (service.StoryService, io.Reader) {
				return &service.MockStoriesService{}, nil
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"message\":\"unexpected end of JSON input\"},\"success\":false}",
		},
		"test add story failure when body is empty": {
			input: func() (service.StoryService, io.Reader) {
				st := contract.AddStoryRequest{Title: "title", Body: ""}
				b, err := json.Marshal(&st)
				require.NoError(t, err)

				return &service.MockStoriesService{}, bytes.NewBuffer(b)
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"message\":\"body cannot be empty\"},\"success\":false}",
		},
		"test add story failure when title is empty": {
			input: func() (service.StoryService, io.Reader) {
				st := contract.AddStoryRequest{Title: "", Body: "test body"}
				b, err := json.Marshal(&st)
				require.NoError(t, err)

				return &service.MockStoriesService{}, bytes.NewBuffer(b)
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"message\":\"title cannot be empty\"},\"success\":false}",
		},
		"test add story failure when svc call fails": {
			input: func() (service.StoryService, io.Reader) {
				ms := &service.MockStoriesService{}
				ms.On("AddStory", mock.AnythingOfType("*model.Story")).Return(liberr.WithArgs(liberr.SeverityError, errors.New("failed to add story")))

				reqSt := contract.AddStoryRequest{Title: "title", Body: "test body"}
				b, err := json.Marshal(&reqSt)
				require.NoError(t, err)

				return ms, bytes.NewBuffer(b)
			},
			expectedCode:   http.StatusInternalServerError,
			expectedResult: "{\"error\":{\"message\":\"internal server error\"},\"success\":false}",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			svc, body := testCase.input()
			testAddStoryHandler(t, testCase.expectedCode, testCase.expectedResult, svc, body)
		})
	}
}

func testAddStoryHandler(t *testing.T, expectedCode int, expectedBody string, svc service.StoryService, body io.Reader) {
	cfg := config.NewConfig("../../../../local.env")

	ah := handler.NewAddHandler(cfg.StoryConfig(), svc)

	r := httptest.NewRequest(http.MethodPost, "/story/add", body)

	w := httptest.NewRecorder()

	mdl.WithError(reporters.NewLogger("dev", "debug"), ah.AddStory)(w, r)

	assert.Equal(t, expectedCode, w.Code)
	assert.Equal(t, expectedBody, w.Body.String())
}
