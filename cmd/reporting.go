package main

import (
	"fmt"
	newrelic "github.com/newrelic/go-agent"
	"github.com/nsnikhil/stories/cmd/config"
	"go.uber.org/zap"
	"gopkg.in/alexcesaro/statsd.v2"
)

var logger *zap.Logger
var nrApp newrelic.Application
var sc *statsd.Client

func initReporters(cfg config.Config) {
	initLogger(cfg.GetEnv())
	initNewRelic(cfg.GetNewRelicConfig())
	initStatsD(cfg.GetStatsDConfig())
}

func initLogger(env string) {
	var err error

	if env == dev {
		logger, err = zap.NewDevelopmentConfig().Build()
	} else {
		logger, err = zap.NewProductionConfig().Build()
	}

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Println(err)
		}
	}()
}

func initNewRelic(nrc config.NewRelicConfig) {
	var err error
	nrApp, err = newrelic.NewApplication(newrelic.NewConfig(nrc.AppName(), nrc.LicenseKey()))
	if err != nil {
		panic(err)
	}
}

func initStatsD(sdc config.StatsDConfig) {
	var err error
	sc, err = statsd.New(statsd.Address(sdc.Address()), statsd.Prefix(sdc.Namespace()))
	if err != nil {
		panic(err)
	}
}
