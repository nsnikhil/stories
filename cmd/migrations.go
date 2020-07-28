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

func Migrate() error {
	newMigrate, err := newMigrate()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return newMigrate.Up()
}

func Rollback() error {
	newMigrate, err := newMigrate()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return newMigrate.Steps(rollBackStep)
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
