package server

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgrpc"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/cmd/config"
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
	deps *deps
	proto.UnimplementedStoriesApiServer
}

func newStoriesServer(deps *deps) *storiesServer {
	return &storiesServer{
		deps: deps,
	}
}

func StartServer(cfg config.Config, logger *zap.Logger, nrApp newrelic.Application, sc *statsd.Client) {
	listener, err := net.Listen("tcp", cfg.GetServerConfig().Address())
	if err != nil {
		logger.Sugar().Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Minute * time.Duration(cfg.GetServerConfig().IdleConnectionTimeoutInMinutes()),
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

	deps := newDeps(getService(cfg, logger), cfg, logger)
	storiesServer := newStoriesServer(deps)
	healthServer := newHealthServer()
	proto.RegisterStoriesApiServer(grpcServer, storiesServer)
	proto.RegisterHealthServer(grpcServer, healthServer)

	logger.Sugar().Infof("listening on %s", cfg.GetServerConfig().Address())
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
