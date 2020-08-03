package server

import (
	"github.com/nsnikhil/stories-proto/proto"
)

type storiesServer struct {
	proto.UnimplementedStoriesApiServer
	deps *serverDeps
}

func newStoriesServer(deps *serverDeps) *storiesServer {
	return &storiesServer{
		deps: deps,
	}
}
