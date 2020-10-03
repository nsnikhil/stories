package stories

import (
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/story/service"
)

//goland:noinspection ALL
type StoriesServer struct {
	proto.UnimplementedStoriesApiServer
	cfg config.StoryConfig
	svc service.StoryService
}

func NewStoriesServer(cfg config.StoryConfig, svc service.StoryService) *StoriesServer {
	return &StoriesServer{
		cfg: cfg,
		svc: svc,
	}
}
