package main

import (
	"github.com/golang-migrate/migrate/v4"
)

func runMigrations() {
	if err := Migrate(); err != nil {
		if err == migrate.ErrNoChange {
			return
		}
		return
	}
}
