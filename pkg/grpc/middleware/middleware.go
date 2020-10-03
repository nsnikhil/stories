package middleware

import (
	"context"
	reporters "github.com/nsnikhil/stories/pkg/reporting"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"strings"
	"time"
)

func WithErrorLogger(lgr *zap.Logger) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		h, err := handler(ctx, req)
		if err != nil {
			lgr.Error(err.Error())
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
