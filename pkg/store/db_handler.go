package store

import (
	"database/sql"
	"github.com/nsnikhil/stories/pkg/config"
	"time"
)

type DBHandler interface {
	GetDB() (*sql.DB, error)
}

type defaultDBHandler struct {
	cfg config.DatabaseConfig
}

func (dbh *defaultDBHandler) GetDB() (*sql.DB, error) {
	db, err := sql.Open(dbh.cfg.DriverName(), dbh.cfg.Source())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(dbh.cfg.MaxOpenConnections())
	db.SetMaxIdleConns(dbh.cfg.IdleConnections())
	db.SetConnMaxLifetime(time.Minute * time.Duration(dbh.cfg.ConnectionMaxLifetime()))

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewDBHandler(cfg config.DatabaseConfig) DBHandler {
	return &defaultDBHandler{
		cfg: cfg,
	}
}
