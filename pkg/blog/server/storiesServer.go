package server

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgrpc"
	"github.com/nsnikhil/stories/pkg/blog/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gopkg.in/alexcesaro/statsd.v2"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type storiesServer struct {
	logger *zap.Logger
}

func newStoriesServer(logger *zap.Logger) *storiesServer {
	return &storiesServer{
		logger: logger,
	}
}

func StartServer(address string, logger *zap.Logger, nrApp newrelic.Application, sc *statsd.Client) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Sugar().Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				withStatsD(sc),
				nrgrpc.UnaryServerInterceptor(nrApp),
				grpc_recovery.UnaryServerInterceptor(),
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				nrgrpc.StreamServerInterceptor(nrApp),
				grpc_recovery.StreamServerInterceptor(),
			),
		),
	)

	storiesServer := newStoriesServer(logger)
	healthServer := newHealthServer()
	proto.RegisterStoriesApiServer(grpcServer, storiesServer)
	proto.RegisterHealthServer(grpcServer, healthServer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			panic(err)
		}
	}()

	waitForShutdown(grpcServer)
}

func waitForShutdown(grpcServer *grpc.Server) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	grpcServer.GracefulStop()
}

func (ss *storiesServer) Ping(ctx context.Context, in *proto.PingRequest) (*proto.PingResponse, error) {
	ss.logger.Info("[storiesServer] [Ping]")
	return &proto.PingResponse{Message: "pong"}, nil
}
