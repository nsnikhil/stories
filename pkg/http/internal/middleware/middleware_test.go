package middleware_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/nsnikhil/stories/pkg/http/internal/middleware"
	"github.com/nsnikhil/stories/pkg/liberr"
	reporters "github.com/nsnikhil/stories/pkg/reporting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWithErrorHandling(t *testing.T) {
	testCases := []struct {
		name           string
		handler        func(resp http.ResponseWriter, req *http.Request) error
		expectedResult string
		expectedCode   int
		expectedLog    string
	}{
		{
			name: "test error middleware with typed error",
			handler: func(resp http.ResponseWriter, req *http.Request) error {
				db := func() error {
					return liberr.WithArgs(
						liberr.Operation("db.insert"),
						liberr.Kind("databaseError"),
						liberr.SeverityError,
						errors.New("insertion failed"),
					)
				}

				svc := func() error {
					return liberr.WithArgs(
						liberr.Operation("svc.addUser"),
						liberr.Kind("dependencyError"),
						liberr.SeverityWarn,
						db(),
					)
				}

				return liberr.WithArgs(
					liberr.Operation("handler.addUser"),
					liberr.ValidationError,
					liberr.SeverityInfo,
					svc(),
				)
			},
			expectedResult: "{\"error\":{\"message\":\"insertion failed\"},\"success\":false}",
			expectedCode:   http.StatusBadRequest,
			expectedLog:    "insertion failed",
		},
		{
			name: "test error middleware with error",
			handler: func(resp http.ResponseWriter, req *http.Request) error {
				return errors.New("some random error")
			},
			expectedResult: "{\"error\":{\"message\":\"internal server error\"},\"success\":false}",
			expectedCode:   http.StatusInternalServerError,
			expectedLog:    "some random error",
		},
		{
			name: "test error middleware with no error",
			handler: func(resp http.ResponseWriter, req *http.Request) error {
				resp.WriteHeader(http.StatusOK)
				_, _ = resp.Write([]byte("success"))
				return nil
			},
			expectedResult: "success",
			expectedCode:   http.StatusOK,
			expectedLog:    "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testWithError(t, testCase.expectedCode, testCase.expectedResult, testCase.expectedLog, testCase.handler)
		})
	}
}

func testWithError(t *testing.T, expectedCode int, expectedBody, expectedLog string, h func(http.ResponseWriter, *http.Request) error) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/random", nil)
	require.NoError(t, err)

	buf := new(bytes.Buffer)

	lgr := reporters.NewLogger("dev", "debug", buf)

	middleware.WithError(lgr, h)(w, r)

	assert.Equal(t, expectedCode, w.Code)
	assert.Equal(t, expectedBody, w.Body.String())

	if len(expectedLog) != 0 {
		assert.True(t, strings.Contains(buf.String(), expectedLog))
	}
}

func TestWithPrometheus(t *testing.T) {
	type prometheusTest struct {
		method   string
		argument []interface{}
	}

	pt := func(method string, args ...interface{}) prometheusTest {
		return prometheusTest{
			method:   method,
			argument: args,
		}
	}

	testCases := []struct {
		name         string
		actualResult func() (*reporters.MockPrometheus, []prometheusTest)
	}{
		{
			name: "test prometheus middleware for success",
			actualResult: func() (*reporters.MockPrometheus, []prometheusTest) {
				w := httptest.NewRecorder()
				r, err := http.NewRequest(http.MethodGet, "/random", nil)
				require.NoError(t, err)

				th := func(resp http.ResponseWriter, req *http.Request) {
					resp.WriteHeader(http.StatusOK)
				}

				mockPrometheus := &reporters.MockPrometheus{}
				mockPrometheus.On("ReportAttempt", "random")
				mockPrometheus.On("ReportSuccess", "random")
				mockPrometheus.On("Observe", "random", mock.Anything)

				middleware.WithPrometheus(mockPrometheus, "random", th)(w, r)

				return mockPrometheus, []prometheusTest{
					pt("ReportAttempt", "random"),
					pt("ReportSuccess", "random"),
					pt("Observe", "random", mock.Anything),
				}
			},
		},
		{
			name: "test prometheus middleware for 400 error",
			actualResult: func() (*reporters.MockPrometheus, []prometheusTest) {
				w := httptest.NewRecorder()
				r, err := http.NewRequest(http.MethodGet, "/random", nil)
				require.NoError(t, err)

				th := func(resp http.ResponseWriter, req *http.Request) {
					resp.WriteHeader(http.StatusBadRequest)
				}

				mockPrometheus := &reporters.MockPrometheus{}
				mockPrometheus.On("ReportAttempt", "random")
				mockPrometheus.On("ReportFailure", "random")
				mockPrometheus.On("Observe", "random", mock.Anything)

				middleware.WithPrometheus(mockPrometheus, "random", th)(w, r)

				return mockPrometheus, []prometheusTest{
					pt("ReportAttempt", "random"),
					pt("ReportFailure", "random"),
					pt("Observe", "random", mock.Anything),
				}
			},
		},
		{
			name: "test statsd middleware for 500 error",
			actualResult: func() (*reporters.MockPrometheus, []prometheusTest) {
				w := httptest.NewRecorder()
				r, err := http.NewRequest(http.MethodGet, "/random", nil)
				require.NoError(t, err)

				th := func(resp http.ResponseWriter, req *http.Request) {
					resp.WriteHeader(http.StatusInternalServerError)
				}

				mockPrometheus := &reporters.MockPrometheus{}
				mockPrometheus.On("ReportAttempt", "random")
				mockPrometheus.On("ReportFailure", "random")
				mockPrometheus.On("Observe", "random", mock.Anything)

				middleware.WithPrometheus(mockPrometheus, "random", th)(w, r)

				return mockPrometheus, []prometheusTest{
					pt("ReportAttempt", "random"),
					pt("ReportFailure", "random"),
					pt("Observe", "random", mock.Anything),
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cl, res := testCase.actualResult()
			for _, r := range res {
				cl.AssertCalled(t, r.method, r.argument...)
			}
		})
	}
}

func TestWithReqRespLog(t *testing.T) {
	type CusReq struct {
		ReqID   string `json:"req_id"`
		ReqData string `json:"req_data"`
	}

	type CusResp struct {
		RespID   string `json:"resp_id"`
		RespData string `json:"resp_data"`
	}

	cReq := CusReq{ReqID: "req-id", ReqData: "req data"}

	b, err := json.Marshal(&cReq)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/random", bytes.NewBuffer(b))
	require.NoError(t, err)

	th := func(resp http.ResponseWriter, req *http.Request) {
		cResp := CusResp{RespID: "resp-id", RespData: "resp data"}

		b, err := json.Marshal(&cResp)
		require.NoError(t, err)

		resp.WriteHeader(http.StatusCreated)
		resp.Header().Set("X-VALUE", "value")
		_, _ = resp.Write(b)
	}

	buf := new(bytes.Buffer)

	lg := reporters.NewLogger("dev", "debug", buf)

	middleware.WithReqRespLog(lg, th)(w, r)

	//assert.True(t, strings.Contains(buf.String(), `{\"req_id\":\"req-id\",\"req_data\":\"req data\"}`))
	//assert.True(t, strings.Contains(buf.String(), `{\"resp_id\":\"resp-id\",\"resp_data\":\"resp data\"}`))
}

func TestWithResponseHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/random", nil)
	require.NoError(t, err)

	th := func(resp http.ResponseWriter, req *http.Request) {}

	middleware.WithResponseHeaders(th)(w, r)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestWithRequestContext(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/random", nil)
	require.NoError(t, err)

	th := func(resp http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "val", req.Context().Value("key"))
	}

	middleware.WithRequestContext(th)(w, r)
}
