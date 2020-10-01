package main

import (
	"github.com/nsnikhil/stories/pkg/app"
	"github.com/nsnikhil/stories/pkg/store"
	"log"
)

const (
	serveCommand    = "serve"
	migrateCommand  = "migrate"
	rollbackCommand = "rollback"
)

func commands() map[string]func() {
	return map[string]func(){
		serveCommand:    app.Start,
		migrateCommand:  store.RunMigrations,
		rollbackCommand: store.RollBackMigrations,
	}
}

func execute(cmd string) {
	run, ok := commands()[cmd]
	if !ok {
		log.Fatal("invalid command")
	}

	run()
}
