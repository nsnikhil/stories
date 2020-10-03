package app

func StartGRPCServer() {
	initGRPCServer().Start()
}

func StartHTTPServer() {
	initHTTPServer().Start()
}
