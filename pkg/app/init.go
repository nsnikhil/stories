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

func initGRPCServer(configFile string) grpcserver.Server {
	cfg, lgr, pr, nr, svc := initCommons(configFile)
	return grpcserver.NewServer(cfg, lgr, nr, pr, svc)
}

func initHTTPServer(configFile string) httpserver.Server {
	cfg, lgr, pr, nr, svc := initCommons(configFile)
	rt := initRouter(cfg.StoryConfig(), lgr, nr, pr, svc)
	return httpserver.NewServer(cfg, lgr, rt)
}

func initCommons(configFile string) (config.Config, *zap.Logger, reporters.Prometheus, *newrelic.Application, service.StoryService) {
	cfg := config.NewConfig(configFile)

	lgr := initLogger(cfg)
	pr := reporters.NewPrometheus()
	nr, err := reporters.NewNewRelicApp(cfg.NewRelicConfig())
	if err != nil {
		log.Fatal(err)
	}

	svc := initService(cfg)

	return cfg, lgr, pr, nr, svc
}

func initRouter(cfg config.StoryConfig, lgr *zap.Logger, newRelic *newrelic.Application, prometheus reporters.Prometheus, svc service.StoryService) http.Handler {
	return router.NewRouter(cfg, lgr, newRelic, prometheus, svc)
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
		getWriters(cfg)...,
	)
}

func getWriters(cfg config.Config) []io.Writer {
	//TODO: MOVE TO CONST
	logSinkMap := map[string]io.Writer{
		"stdout": os.Stdout,
		"file":   reporters.NewExternalLogFile(cfg.LogFileConfig()),
	}

	var writers []io.Writer
	for _, sink := range cfg.LogConfig().Sinks() {
		w, ok := logSinkMap[sink]
		if ok {
			writers = append(writers, w)
		}
	}

	return writers
}
