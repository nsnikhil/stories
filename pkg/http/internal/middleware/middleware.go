package middleware

import (
	"context"
	"fmt"
	"github.com/nsnikhil/stories/pkg/http/internal/resperr"
	"github.com/nsnikhil/stories/pkg/http/internal/util"
	"github.com/nsnikhil/stories/pkg/liberr"
	reporters "github.com/nsnikhil/stories/pkg/reporting"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func WithError(handler func(resp http.ResponseWriter, req *http.Request) error) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {

		err := handler(resp, req)
		if err == nil {
			return
		}

		t, ok := err.(*liberr.Error)
		if ok {
			fmt.Println(t.Details())
		} else {
			fmt.Println(err)
		}

		util.WriteFailureResponse(resperr.MapError(err), resp)
	}
}

func WithReqRespLog(lgr *zap.Logger, handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		cr := util.NewCopyWriter(resp)

		handler(cr, req)

		//b, _ := cr.Body()

		//lgr.Sugar().Debug(req)
		//lgr.Sugar().Debug(string(b))
	}
}

func WithResponseHeaders(handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "application/json")
		handler(resp, req)
	}
}

func WithRequestContext(handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		// TODO: CHANGE TEMP VALUE
		ctx := context.WithValue(req.Context(), "key", "val")
		handler(resp, req.WithContext(ctx))
	}
}

func WithPrometheus(prometheus reporters.Prometheus, api string, handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		// TODO CHANGE THIS
		hasError := func(code int) bool {
			return code >= 400 && code <= 600
		}

		start := time.Now()
		prometheus.ReportAttempt(api)

		cr := util.NewCopyWriter(resp)

		handler(cr, req)
		if hasError(cr.Code()) {
			duration := time.Since(start)
			prometheus.Observe(api, duration.Seconds())
			prometheus.ReportFailure(api)
			return
		}

		duration := time.Since(start)
		prometheus.Observe(api, duration.Seconds())

		prometheus.ReportSuccess(api)
	}
}
