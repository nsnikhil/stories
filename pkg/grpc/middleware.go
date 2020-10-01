package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"gopkg.in/alexcesaro/statsd.v2"
	"strings"
)

const (
	attempt = "attempt"
	success = "success"
	failure = "failure"
)

func withStatsD(sc *statsd.Client) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		tm := sc.NewTiming()
		inc := func(a, m string, sc *statsd.Client) {
			sc.Increment(fmt.Sprintf("%s.%s.counter", a, m))
		}

		method := func(fm string) string {
			r := strings.Split(fm, "/")
			return strings.ToLower(r[len(r)-1])
		}

		defer sc.Flush()
		m := method(info.FullMethod)

		inc(m, attempt, sc)

		h, err := handler(ctx, req)

		if err != nil {
			inc(m, failure, sc)
		} else {
			inc(m, success, sc)
		}

		tm.Send(fmt.Sprintf("%s.%s.counter", m, "time"))
		return h, err
	}
}
