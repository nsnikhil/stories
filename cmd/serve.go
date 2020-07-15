package main

import (
	"github.com/nsnikhil/stories/pkg/blog/server"
)

func init() {
	initConfigs()
	initReporters()
}

func startServer() {
	server.StartServer(cfg.sc.address(), logger, nrApp, sc)
}
