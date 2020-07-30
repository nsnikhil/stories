package server

import (
	"github.com/nsnikhil/stories-proto/proto"
)

type storiesServer struct {
	deps *serverDeps
	proto.UnimplementedStoriesApiServer
}

func newStoriesServer(deps *serverDeps) *storiesServer {
	return &storiesServer{
		deps: deps,
	}
}
