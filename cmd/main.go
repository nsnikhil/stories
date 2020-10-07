package main

import (
	"flag"
	"fmt"
	"log"
)

const (
	configFileKey     = "configFile"
	defaultConfigFile = "local.env"
	configFileUsage   = ""
)

func main() {
	var configFile string
	flag.StringVar(&configFile, configFileKey, defaultConfigFile, configFileUsage)
	flag.Parse()

	if len(configFile) == 0 {
		log.Fatal("please provide config file")
	}

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("please provide a command")
	}

	fmt.Println(flag.Args())

	execute(flag.Args()[0], configFile)
}
