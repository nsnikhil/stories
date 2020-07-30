package server

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgrpc"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/cmd/config"
	"github.com/nsnikhil/stories/pkg/blog/service"
	"github.com/nsnikhil/stories/pkg/blog/store"
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

type Server interface {
	Start()
}

type AppServer struct {
	cfg      config.Config
	logger   *zap.Logger
	newRelic newrelic.Application
	statsd   *statsd.Client
}

func NewAppServer(cfg config.Config, logger *zap.Logger, newRelic newrelic.Application, statsd *statsd.Client) Server {
	return &AppServer{
		cfg:      cfg,
		logger:   logger,
		newRelic: newRelic,
		statsd:   statsd,
	}
}

func (as *AppServer) Start() {
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: time.Minute * time.Duration(as.cfg.GetServerConfig().IdleConnectionTimeoutInMinutes())}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(withStatsD(as.statsd), nrgrpc.UnaryServerInterceptor(as.newRelic), grpc_recovery.UnaryServerInterceptor())),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(nrgrpc.StreamServerInterceptor(as.newRelic), grpc_recovery.StreamServerInterceptor())),
	)

	svc := getService(as.cfg, as.logger)
	serverDeps := newDeps(svc, as.cfg, as.logger)
	storiesServer := newStoriesServer(serverDeps)

	healthServer := newHealthServer()

	proto.RegisterStoriesApiServer(grpcServer, storiesServer)
	proto.RegisterHealthServer(grpcServer, healthServer)

	listener, err := net.Listen("tcp", as.cfg.GetServerConfig().Address())
	if err != nil {
		as.logger.Sugar().Fatalf("failed to listen: %v", err)
	}

	as.logger.Sugar().Infof("listening on %s", as.cfg.GetServerConfig().Address())
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

type serverDeps struct {
	cfg    config.Config
	logger *zap.Logger
	svc    *service.Service
}

func newDeps(svc *service.Service, cfg config.Config, logger *zap.Logger) *serverDeps {
	return &serverDeps{
		cfg:    cfg,
		svc:    svc,
		logger: logger,
	}
}

func getService(cfg config.Config, lgr *zap.Logger) *service.Service {
	handler := store.NewDBHandler(cfg.GetDatabaseConfig(), lgr)
	db, err := handler.GetDB()
	if err != nil {
		panic(err)
	}

	trie := store.NewCharacterTrie()

	return service.NewService(
		service.NewDefaultStoriesService(store.NewStore(store.NewDefaultStoriesStore(db, lgr), store.NewTrieStoriesCache(trie, cfg.GetBlogConfig(), lgr)), lgr),
	)
}
