package store

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nsnikhil/stories/cmd/config"
	"go.uber.org/zap"
)

const databaseDriverName = "postgres"

type DBHandler interface {
	GetDB() (*gorm.DB, error)
}

type DefaultDBHandler struct {
	config *config.DatabaseConfig
	logger *zap.Logger
}

func NewDBHandler(config *config.DatabaseConfig, logger *zap.Logger) DBHandler {
	return &DefaultDBHandler{
		config: config,
		logger: logger,
	}
}

func (dbh *DefaultDBHandler) GetDB() (*gorm.DB, error) {
	db, err := gorm.Open(databaseDriverName, dbh.config.Source())
	if err != nil {
		dbh.logger.Error(err.Error(), zap.String("method", "GetDB"), zap.String("call", "gorm.Open"))
		return nil, err
	}

	return db, nil
}
