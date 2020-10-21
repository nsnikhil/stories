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

	rf := func(method, path string) *http.Request {
		req, err := http.NewRequest(method, path, nil)
		require.NoError(t, err)

		return req
	}

	testCases := map[string]struct {
		name    string
		request *http.Request
	}{
		"test ping route": {
			request: rf(http.MethodGet, "/ping"),
		},
		"test add story route": {
			request: rf(http.MethodPost, "/story/add"),
		},
		"test get story route": {
			request: rf(http.MethodPost, "/story/get"),
		},
		"test delete story route": {
			request: rf(http.MethodPost, "/story/delete"),
		},
		"test most viewed stories route": {
			request: rf(http.MethodPost, "/story/most-viewed"),
		},
		"test top rated stories route": {
			request: rf(http.MethodPost, "/story/top-rated"),
		},
		"test search stories route": {
			request: rf(http.MethodPost, "/story/search"),
		},
		"test update story route": {
			request: rf(http.MethodPost, "/story/update"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()

			r.ServeHTTP(w, testCase.request)

			assert.NotEqual(t, http.StatusNotFound, w.Code)
		})
	}
}
