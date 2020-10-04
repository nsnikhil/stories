package stories

import (
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/story/service"
)

type Server struct {
	proto.UnimplementedStoriesApiServer
	cfg config.StoryConfig
	svc service.StoryService
}

func NewStoriesServer(cfg config.StoryConfig, svc service.StoryService) *Server {
	return &Server{
		cfg: cfg,
		svc: svc,
	}
}
