package grpc

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"go.uber.org/zap"
)

func (ss *storiesServer) SearchStories(ctx context.Context, req *proto.SearchStoriesRequest) (*proto.SearchStoriesResponse, error) {
	stories, err := ss.deps.svc.GetStoriesService().SearchStories(req.GetQuery())
	if err != nil {
		return nil, logAndGetError(err, "SearchStories", "SearchStories", ss.deps.logger)
	}

	sz := len(stories)
	resp := make([]*proto.Story, sz)

	for i := 0; i < sz; i++ {
		resp[i] = toProtoStory(&stories[i])
	}

	ss.deps.logger.Info("search successfully", zap.String("method", "SearchStories"))
	return &proto.SearchStoriesResponse{Stories: resp}, nil
}
