package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/nsnikhil/stories/cmd/config"
	"github.com/nsnikhil/stories/pkg/blog/store"
	"go.uber.org/zap"
	"path/filepath"
	"strings"
)

const (
	migrationPath = "./pkg/blog/store/migrations"
	rollBackStep  = -1
	cutSet        = "file://"
	databaseName  = "postgres"
)

func runMigrations() {
	newMigrate, err := newMigrate()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := newMigrate.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return
		}
		fmt.Println(err)
		return
	}
}

func rollBackMigrations() {
	newMigrate, err := newMigrate()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := newMigrate.Steps(rollBackStep); err != nil {
		if err == migrate.ErrNoChange {
			return
		}
	}
}

func newMigrate() (*migrate.Migrate, error) {
	cfg := config.LoadConfigs()

	dbHandler := store.NewDBHandler(cfg.GetDatabaseConfig(), zap.NewExample())

	db, err := dbHandler.GetDB()
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db.DB(), &postgres.Config{})
	if err != nil {
		return nil, err
	}

	sourcePath, err := getSourcePath(migrationPath)
	if err != nil {
		return nil, err
	}

	return migrate.NewWithDatabaseInstance(sourcePath, databaseName, driver)
}

func getSourcePath(directory string) (string, error) {
	directory = strings.TrimLeft(directory, cutSet)
	absPath, err := filepath.Abs(directory)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", cutSet, absPath), nil
}
