package server

import (
	"github.com/nsnikhil/stories/cmd/config"
	"github.com/nsnikhil/stories/pkg/blog/service"
	"github.com/nsnikhil/stories/pkg/blog/store"
	"go.uber.org/zap"
)

type deps struct {
	cfg    config.Config
	logger *zap.Logger
	svc    *service.Service
}

func newDeps(svc *service.Service, cfg config.Config, logger *zap.Logger) *deps {
	return &deps{
		cfg:    cfg,
		svc:    svc,
		logger: logger,
	}
}

func getService(cfg config.Config, logger *zap.Logger) *service.Service {
	return service.NewService(service.NewDefaultStoriesService(getStore(cfg, logger), logger))
}

func getStore(cfg config.Config, logger *zap.Logger) *store.Store {
	return store.NewStore(getDB(cfg, logger), getCache(logger))
}

func getCache(logger *zap.Logger) store.StoriesCache {
	tr := store.NewCharacterTrie()
	return store.NewTrieStoriesCache(tr, logger)
}

func getDB(cfg config.Config, logger *zap.Logger) store.StoriesStore {
	handler := store.NewDBHandler(cfg.GetDatabaseConfig(), logger)
	db, err := handler.GetDB()
	if err != nil {
		panic(err)
	}

	return store.NewDefaultStoriesStore(db, logger)
}
