package main

import (
	"github.com/golang-migrate/migrate/v4"
)

func rollBackMigrations() {
	if err := Rollback(); err != nil {
		if err == migrate.ErrNoChange {
			return
		}
	}
}
