package app

import (
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/grpc/server"
	reporters "github.com/nsnikhil/stories/pkg/reporting"
)

func Start() {
	cfg := config.NewConfig()

	lgr := initLogger(cfg)
	pr := reporters.NewPrometheus()
	nr := reporters.NewNewRelicApp(cfg.NewRelicConfig())

	svc := initService(cfg)

	server.NewServer(cfg, lgr, nr, pr, svc).Start()
}
