package main

import (
	"fmt"
	"os"
)

func start(command string) {
	commands := map[string]func(){
		"serve": startServer,
	}

	if commands[command] == nil {
		panic(fmt.Errorf("invalid command %s", command))
	}

	newCommand(command, commands[command]).execute()
}

func main() {
	args := os.Args
	if len(args) == 0 {
		panic("please provide an valid command")
	}
	start(args[1])
}
