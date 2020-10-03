package app

import (
	"github.com/nsnikhil/stories/pkg/config"
	reporters "github.com/nsnikhil/stories/pkg/reporting"
	"github.com/nsnikhil/stories/pkg/store"
	"github.com/nsnikhil/stories/pkg/story/service"
	"go.uber.org/zap"
	"io"
	"log"
	"os"
)

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
