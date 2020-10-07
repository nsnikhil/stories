package app

func StartGRPCServer(configFile string) {
	initGRPCServer(configFile).Start()
}

func StartHTTPServer(configFile string) {
	initHTTPServer(configFile).Start()
}
