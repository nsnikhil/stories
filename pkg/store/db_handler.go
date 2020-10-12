package store

import (
	"database/sql"
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/liberr"
	"time"
)

const ()

type DBHandler interface {
	GetDB() (*sql.DB, error)
}

type defaultDBHandler struct {
	cfg config.DatabaseConfig
}

func (dbh *defaultDBHandler) GetDB() (*sql.DB, error) {
	db, err := sql.Open(dbh.cfg.DriverName(), dbh.cfg.Source())
	if err != nil {
		return nil, liberr.WithArgs(liberr.Operation("DBHandler.GetDB.sql.Open"), liberr.SeverityError, err)
	}

	db.SetMaxOpenConns(dbh.cfg.MaxOpenConnections())
	db.SetMaxIdleConns(dbh.cfg.IdleConnections())
	db.SetConnMaxLifetime(time.Minute * time.Duration(dbh.cfg.ConnectionMaxLifetime()))

	if err := db.Ping(); err != nil {
		return nil, liberr.WithArgs(liberr.Operation("DBHandler.GetDB.db.Ping"), liberr.SeverityError, err)
	}

	return db, nil
}

func NewDBHandler(cfg config.DatabaseConfig) DBHandler {
	return &defaultDBHandler{
		cfg: cfg,
	}
}
