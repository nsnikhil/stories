package main

import (
	"github.com/nsnikhil/stories/cmd/config"
	"github.com/nsnikhil/stories/pkg/blog/server"
)

func startServer() {
	cfg := config.LoadConfigs()
	initReporters(cfg)
	server.StartServer(cfg.GetServerConfig(), logger, nrApp, sc)
}
