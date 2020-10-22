package middleware

import (
	"context"
	"github.com/nsnikhil/stories/pkg/liberr"
	reporters "github.com/nsnikhil/stories/pkg/reporting"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"strings"
	"time"
)

//TODO: ADD MASKING BEFORE LOGGING REQ AND RESP
func WithReqRespLogger(lgr *zap.Logger) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		//lgr.Sugar().Debug(req)

		h, err := handler(ctx, req)
		if err != nil {
			//lgr.Sugar().Debug(err)
			return h, err
		}

		//lgr.Sugar().Debug(h)
		return h, err
	}
}

func WithErrorLogger(lgr *zap.Logger) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		h, err := handler(ctx, req)
		if err != nil {
			t, ok := err.(*liberr.Error)
			if ok {
				lgr.Error(t.EncodedStack())
			} else {
				lgr.Error(err.Error())
			}
		}

		return h, err
	}
}

func WithPrometheus(pr reporters.Prometheus) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		methods := strings.Split(info.FullMethod, "/")
		api := strings.ToLower(methods[len(methods)-1])

		now := time.Now()
		pr.ReportAttempt(api)

		h, err := handler(ctx, req)

		if err != nil {
			pr.ReportFailure(api)
		} else {
			pr.ReportSuccess(api)
		}

		duration := time.Since(now)
		pr.Observe(api, duration.Seconds())

		return h, err
	}
}
