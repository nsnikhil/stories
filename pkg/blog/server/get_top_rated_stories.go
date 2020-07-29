package server

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"go.uber.org/zap"
)

func (ss *storiesServer) GetTopRatedStories(ctx context.Context, req *proto.TopRatedStoriesRequest) (*proto.TopRatedStoriesResponse, error) {
	stories, err := ss.deps.svc.GetStoriesService().GetTopRatedStories(int(req.GetOffset()), int(req.GetLimit()))
	if err != nil {
		return nil, logAndGetError(err, "GetTopRatedStories", "GetTopRatedStories", ss.deps.logger)
	}

	sz := len(stories)
	resp := make([]*proto.Story, sz)

	for i := 0; i < sz; i++ {
		resp[i] = toProtoStory(&stories[i])
	}

	ss.deps.logger.Info("get top rated stories successfully", zap.String("method", "GetTopRatedStories"))
	return &proto.TopRatedStoriesResponse{Stories: resp}, nil
}
