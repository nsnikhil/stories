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
)

func TestAddStory(t *testing.T) {
	cfg := config.NewConfig()

	testCases := []struct {
		name           string
		actualResult   func() (string, int)
		expectedResult string
		expectedCode   int
	}{
		{
			name: "test add story success",
			actualResult: func() (string, int) {
				reqSt := contract.AddStoryRequest{
					Title: "title",
					Body:  "test body",
				}

				st, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("AddStory", st).Return(nil)

				ah := handler.NewAddHandler(cfg.StoryConfig(), ms)

				b, err := json.Marshal(&reqSt)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/story/add", bytes.NewBuffer(b))

				mdl.WithError(ah.Add)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusCreated,
			expectedResult: "{\"data\":{\"success\":true},\"success\":true}",
		},
		{
			name: "test add story fails when body is nil",
			actualResult: func() (string, int) {
				ah := handler.NewAddHandler(cfg.StoryConfig(), &service.MockStoriesService{})

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/story/add", nil)

				mdl.WithError(ah.Add)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"code\":\"usx0001\",\"message\":\"unexpected end of JSON input\"},\"success\":false}",
		},
		{
			name: "test add story failure when title is empty",
			actualResult: func() (string, int) {
				st := contract.AddStoryRequest{
					Title: "",
					Body:  "test body",
				}

				ah := handler.NewAddHandler(cfg.StoryConfig(), &service.MockStoriesService{})

				b, err := json.Marshal(&st)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/story/add", bytes.NewBuffer(b))

				mdl.WithError(ah.Add)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"code\":\"usx0001\",\"message\":\"title cannot be empty\"},\"success\":false}",
		},
		{
			name: "test add story failure when body is empty",
			actualResult: func() (string, int) {
				st := contract.AddStoryRequest{
					Title: "title",
					Body:  "",
				}

				ah := handler.NewAddHandler(cfg.StoryConfig(), &service.MockStoriesService{})

				b, err := json.Marshal(&st)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/story/add", bytes.NewBuffer(b))

				mdl.WithError(ah.Add)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: "{\"error\":{\"code\":\"usx0001\",\"message\":\"body cannot be empty\"},\"success\":false}",
		},
		{
			name: "test add story failure when svc call fails",
			actualResult: func() (string, int) {
				reqSt := contract.AddStoryRequest{
					Title: "title",
					Body:  "test body",
				}

				st, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "test body").
					Build()

				require.NoError(t, err)

				ms := &service.MockStoriesService{}
				ms.On("AddStory", st).Return(errors.New("failed to add story"))

				ah := handler.NewAddHandler(cfg.StoryConfig(), ms)

				b, err := json.Marshal(&reqSt)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/story/add", bytes.NewBuffer(b))

				mdl.WithError(ah.Add)(w, r)

				return w.Body.String(), w.Code
			},
			expectedCode:   http.StatusInternalServerError,
			expectedResult: "{\"error\":{\"code\":\"usx0010\",\"message\":\"failed to add story\"},\"success\":false}",
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
