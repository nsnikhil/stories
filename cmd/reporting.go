package main

import (
	"fmt"
	newrelic "github.com/newrelic/go-agent"
	"go.uber.org/zap"
	"gopkg.in/alexcesaro/statsd.v2"
)

var logger *zap.Logger
var nrApp newrelic.Application
var sc *statsd.Client

func initReporters() {
	initLogger()
	initNewRelic()
	initStatsD()
}

func initLogger() {
	var err error

	if cfg.env == dev {
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

func initNewRelic() {
	var err error
	nrApp, err = newrelic.NewApplication(newrelic.NewConfig(cfg.nr.appName, cfg.nr.licenseKey))
	if err != nil {
		panic(err)
	}
}

func initStatsD() {
	var err error
	sc, err = statsd.New(statsd.Address(cfg.sdc.address()), statsd.Prefix(cfg.sdc.namespace))
	if err != nil {
		panic(err)
	}
}
