package server

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgrpc"
	"github.com/nsnikhil/stories/cmd/config"
	"github.com/nsnikhil/stories/pkg/blog/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"gopkg.in/alexcesaro/statsd.v2"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type storiesServer struct {
	logger *zap.Logger
}

func newStoriesServer(logger *zap.Logger) *storiesServer {
	return &storiesServer{
		logger: logger,
	}
}

func StartServer(cfg config.ServerConfig, logger *zap.Logger, nrApp newrelic.Application, sc *statsd.Client) {
	listener, err := net.Listen("tcp", cfg.Address())
	if err != nil {
		logger.Sugar().Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Minute * time.Duration(cfg.IdleConnectionTimeoutInMinutes()),
		}),
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

	logger.Sugar().Infof("listening on %s", cfg.Address())
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
