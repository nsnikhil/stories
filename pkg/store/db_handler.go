package store

import (
	"database/sql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nsnikhil/stories/pkg/config"
	"go.uber.org/zap"
	"time"
)

type DBHandler interface {
	GetDB() (*sql.DB, error)
}

type defaultDBHandler struct {
	cfg config.DatabaseConfig
	lgr *zap.Logger
}

func (dbh *defaultDBHandler) GetDB() (*sql.DB, error) {
	db, err := sql.Open(dbh.cfg.DriverName(), dbh.cfg.Source())
	if err != nil {
		dbh.lgr.Error(err.Error())
		return nil, err
	}

	db.SetMaxOpenConns(dbh.cfg.MaxOpenConnections())
	db.SetMaxIdleConns(dbh.cfg.IdleConnections())
	db.SetConnMaxLifetime(time.Minute * time.Duration(dbh.cfg.ConnectionMaxLifetime()))

	if err := db.Ping(); err != nil {
		dbh.lgr.Error(err.Error())
		return nil, err
	}

	return db, nil
}

func NewDBHandler(cfg config.DatabaseConfig, lgr *zap.Logger) DBHandler {
	return &defaultDBHandler{
		cfg: cfg,
		lgr: lgr,
	}
}
