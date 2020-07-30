package server

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"go.uber.org/zap"
)

func (ss *storiesServer) GetStory(ctx context.Context, req *proto.GetStoryRequest) (*proto.GetStoryResponse, error) {
	st, err := ss.deps.svc.GetStoriesService().GetStory(req.GetStoryID())
	if err != nil {
		return nil, logAndGetError(err, "GetStory", "GetStory", ss.deps.logger)
	}

	ss.deps.logger.Info("get story success", zap.String("method", "GetStory"))
	return &proto.GetStoryResponse{Story: toProtoStory(st)}, nil
}
