package app

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/nsnikhil/stories/pkg/config"
	grpcserver "github.com/nsnikhil/stories/pkg/grpc/server"
	"github.com/nsnikhil/stories/pkg/http/router"
	httpserver "github.com/nsnikhil/stories/pkg/http/server"
	reporters "github.com/nsnikhil/stories/pkg/reporting"
	"github.com/nsnikhil/stories/pkg/store"
	"github.com/nsnikhil/stories/pkg/story/service"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"os"
)

//TODO: COMBINE INIT GRPC AND INIT HTTP SERVER
func initGRPCServer() grpcserver.Server {
	cfg := config.NewConfig()

	lgr := initLogger(cfg)
	pr := reporters.NewPrometheus()
	nr := reporters.NewNewRelicApp(cfg.NewRelicConfig())

	svc := initService(cfg)

	return grpcserver.NewServer(cfg, lgr, nr, pr, svc)
}

//TODO: COMBINE INIT GRPC AND INIT HTTP SERVER
func initHTTPServer() httpserver.Server {
	cfg := config.NewConfig()

	lgr := initLogger(cfg)
	pr := reporters.NewPrometheus()
	nr := reporters.NewNewRelicApp(cfg.NewRelicConfig())

	svc := initService(cfg)

	rt := initRouter(cfg.StoryConfig(), lgr, nr, pr, svc)

	return httpserver.NewServer(cfg, lgr, rt)
}

func initRouter(cfg config.StoryConfig, lgr *zap.Logger, newRelic *newrelic.Application, prometheus reporters.Prometheus, svc service.StoryService) http.Handler {
	return router.NewRouter(
		cfg,
		lgr,
		newRelic,
		prometheus,
		svc,
	)
}

func initService(cfg config.Config) service.StoryService {
	str := initStore(cfg.DatabaseConfig())
	return service.NewStoriesService(str)
}

func initStore(cfg config.DatabaseConfig) store.StoriesStore {
	dbh := store.NewDBHandler(cfg)

	db, err := dbh.GetDB()
	if err != nil {
		log.Fatal(dbh)
	}

	return store.NewStoriesStore(db)
}

func initLogger(cfg config.Config) *zap.Logger {
	return reporters.NewLogger(
		cfg.Env(),
		cfg.LogConfig().Level(),
		getWriters(cfg.LogFileConfig())...,
	)
}

func getWriters(cfg config.LogFileConfig) []io.Writer {
	return []io.Writer{
		os.Stdout,
		reporters.NewExternalLogFile(cfg),
	}
}
