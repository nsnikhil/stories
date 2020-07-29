package server

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"go.uber.org/zap"
)

func (ss *storiesServer) GetMostViewedStories(ctx context.Context, req *proto.MostViewedStoriesRequest) (*proto.MostViewedStoriesResponse, error) {
	stories, err := ss.deps.svc.GetStoriesService().GetMostViewsStories(int(req.GetOffset()), int(req.GetLimit()))
	if err != nil {
		return nil, logAndGetError(err, "GetMostViewedStories", "GetMostViewsStories", ss.deps.logger)
	}

	sz := len(stories)
	resp := make([]*proto.Story, sz)

	for i := 0; i < sz; i++ {
		resp[i] = toProtoStory(&stories[i])
	}

	ss.deps.logger.Info("get most viewed stories successfully", zap.String("method", "GetMostViewedStories"))
	return &proto.MostViewedStoriesResponse{Stories: resp}, nil
}
