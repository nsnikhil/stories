package app

import (
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/grpc"
	"github.com/nsnikhil/stories/pkg/reporting"
)

func Start() {
	cfg := config.NewConfig()
	reporting.initReporters(cfg)

	grpc.NewAppServer(cfg, reporting.logger, reporting.nrApp, reporting.sc).Start()
}
