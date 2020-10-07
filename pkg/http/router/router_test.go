package router_test

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/http/router"
	reporters "github.com/nsnikhil/stories/pkg/reporting"
	"github.com/nsnikhil/stories/pkg/story/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	cfg := config.NewConfig("../../../local.env")

	r := router.NewRouter(
		cfg.StoryConfig(),
		zap.NewNop(),
		&newrelic.Application{},
		&reporters.MockPrometheus{},
		&service.MockStoriesService{},
	)

	testCases := []struct {
		name         string
		actualResult func() int
	}{
		{
			name: "test ping route",
			actualResult: func() int {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodGet, "/ping", nil)
				require.NoError(t, err)

				r.ServeHTTP(resp, req)

				return resp.Code
			},
		},
		{
			name: "test add story route",
			actualResult: func() int {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodPost, "/story/add", nil)
				require.NoError(t, err)

				r.ServeHTTP(resp, req)

				return resp.Code
			},
		},
		{
			name: "test get story route",
			actualResult: func() int {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodPost, "/story/get", nil)
				require.NoError(t, err)

				r.ServeHTTP(resp, req)

				return resp.Code
			},
		},
		{
			name: "test delete story route",
			actualResult: func() int {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodPost, "/story/delete", nil)
				require.NoError(t, err)

				r.ServeHTTP(resp, req)

				return resp.Code
			},
		},
		{
			name: "test most viewed stories route",
			actualResult: func() int {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodPost, "/story/most-viewed", nil)
				require.NoError(t, err)

				r.ServeHTTP(resp, req)

				return resp.Code
			},
		},
		{
			name: "test top rated stories route",
			actualResult: func() int {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodPost, "/story/top-rated", nil)
				require.NoError(t, err)

				r.ServeHTTP(resp, req)

				return resp.Code
			},
		},
		{
			name: "test search stories route",
			actualResult: func() int {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodPost, "/story/search", nil)
				require.NoError(t, err)

				r.ServeHTTP(resp, req)

				return resp.Code
			},
		},
		{
			name: "test update story route",
			actualResult: func() int {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodPost, "/story/update", nil)
				require.NoError(t, err)

				r.ServeHTTP(resp, req)

				return resp.Code
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.NotEqual(t, http.StatusNotFound, testCase.actualResult())
		})
	}
}
