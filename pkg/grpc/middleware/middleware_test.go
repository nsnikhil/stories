package middleware_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/nsnikhil/stories/pkg/grpc/middleware"
	"github.com/nsnikhil/stories/pkg/liberr"
	reporters "github.com/nsnikhil/stories/pkg/reporting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"strings"
	"testing"
)

func TestWithReqRespLogger(t *testing.T) {
	testCases := []struct {
		name            string
		actualResult    func() string
		expectedResults []string
	}{
		{
			name: "test req resp logger for success",
			actualResult: func() string {
				buf := new(bytes.Buffer)
				lgr := reporters.NewLogger("dev", "debug", buf)

				f := middleware.WithReqRespLogger(lgr)

				handler := func(ctx context.Context, req interface{}) (interface{}, error) {
					return "response", nil
				}

				req := "request"

				_, err := f(context.Background(), req, &grpc.UnaryServerInfo{}, handler)
				require.NoError(t, err)

				return buf.String()
			},
			expectedResults: []string{"request", "response"},
		},
		{
			name: "test req resp logger for failure",
			actualResult: func() string {
				buf := new(bytes.Buffer)
				lgr := reporters.NewLogger("dev", "debug", buf)

				f := middleware.WithReqRespLogger(lgr)

				handler := func(ctx context.Context, req interface{}) (interface{}, error) {
					return "", errors.New("some error")
				}

				req := "request"

				_, err := f(context.Background(), req, &grpc.UnaryServerInfo{}, handler)
				require.Error(t, err)

				return buf.String()
			},
			expectedResults: []string{"request", "some error"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			//TODO: ADD TESTS WHEN MASKING IS COMPLETE
			//res := testCase.actualResult()
			//for _, exp := range testCase.expectedResults {
			//	assert.True(t, strings.Contains(res, exp))
			//}
		})
	}
}

func TestWithErrorLogger(t *testing.T) {
	buf := new(bytes.Buffer)

	lgr := reporters.NewLogger("dev", "debug", buf)

	f := middleware.WithErrorLogger(lgr)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
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

		return nil, liberr.WithArgs(
			liberr.Operation("handler.addUser"),
			liberr.InternalError,
			liberr.SeverityInfo,
			svc(),
		)
	}

	req := "request"

	_, err := f(context.Background(), req, &grpc.UnaryServerInfo{}, handler)
	require.Error(t, err)

	assert.True(t, strings.Contains(buf.String(), "insertion failed"))
}

func TestWithPrometheus(t *testing.T) {
	testCases := []struct {
		name     string
		testFunc func()
	}{
		{
			name: "test prometheus for success call",
			testFunc: func() {
				pr := &reporters.MockPrometheus{}
				pr.On("ReportAttempt", "ping")
				pr.On("ReportSuccess", "ping")
				pr.On("Observe", "ping", mock.Anything)

				f := middleware.WithPrometheus(pr)

				handler := func(ctx context.Context, req interface{}) (interface{}, error) {
					return "response", nil
				}

				req := "request"

				_, err := f(context.Background(), req, &grpc.UnaryServerInfo{FullMethod: "/ping"}, handler)
				require.NoError(t, err)
			},
		},
		{
			name: "test prometheus for failure call",
			testFunc: func() {
				pr := &reporters.MockPrometheus{}
				pr.On("ReportAttempt", "ping")
				pr.On("ReportFailure", "ping")
				pr.On("Observe", "ping", mock.Anything)

				f := middleware.WithPrometheus(pr)

				handler := func(ctx context.Context, req interface{}) (interface{}, error) {
					return nil, errors.New("some error")
				}

				req := "request"

				_, err := f(context.Background(), req, &grpc.UnaryServerInfo{FullMethod: "/ping"}, handler)
				require.Error(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.testFunc()
		})
	}
}
