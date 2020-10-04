package server

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/grpc/middleware"
	"github.com/nsnikhil/stories/pkg/grpc/server/health"
	"github.com/nsnikhil/stories/pkg/grpc/server/stories"
	reporters "github.com/nsnikhil/stories/pkg/reporting"
	"github.com/nsnikhil/stories/pkg/story/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server interface {
	Start()
}

type appServer struct {
	cfg config.Config

	lgr *zap.Logger
	pr  reporters.Prometheus
	nr  *newrelic.Application

	svc service.StoryService
}

func NewServer(cfg config.Config, logger *zap.Logger, nr *newrelic.Application, pr reporters.Prometheus, svc service.StoryService) Server {
	return &appServer{
		cfg: cfg,
		lgr: logger,
		pr:  pr,
		nr:  nr,
		svc: svc,
	}
}

func (as *appServer) Start() {
	grpcServer := newGrpcServer(as)

	storiesServer := stories.NewStoriesServer(as.cfg.StoryConfig(), as.svc)
	healthServer := health.NewHealthServer()

	proto.RegisterStoriesApiServer(grpcServer, storiesServer)
	proto.RegisterHealthServer(grpcServer, healthServer)

	setUpPrometheus(as.cfg.GRPCServerConfig(), grpcServer)

	listener, err := net.Listen("tcp", as.cfg.GRPCServerConfig().Address())
	if err != nil {
		as.lgr.Sugar().Fatalf("failed to listen: %v", err)
	}

	as.lgr.Sugar().Infof("listening on %s", as.cfg.GRPCServerConfig().Address())
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	waitForShutdown(grpcServer)
}

func setUpPrometheus(cfg config.GRPCServerConfig, gs *grpc.Server) {
	grpc_prometheus.Register(gs)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		if err := http.ListenAndServe(cfg.PrometheusHTTPAddress(), nil); err != nil {
			fmt.Println(err.Error())
		}
	}()
}

func newGrpcServer(as *appServer) *grpc.Server {
	return grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionIdle: time.Minute * time.Duration(as.cfg.GRPCServerConfig().IdleConnectionTimeoutInMinutes()),
			},
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				middleware.WithReqRespLogger(as.lgr),
				middleware.WithPrometheus(as.pr),
				middleware.WithErrorLogger(as.lgr),
				grpc_recovery.UnaryServerInterceptor(),
				grpc_prometheus.UnaryServerInterceptor,
				nrgrpc.UnaryServerInterceptor(as.nr),
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_recovery.StreamServerInterceptor(),
				grpc_prometheus.StreamServerInterceptor,
				nrgrpc.StreamServerInterceptor(as.nr),
			),
		),
	)
}

func waitForShutdown(grpcServer *grpc.Server) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	grpcServer.GracefulStop()
}
