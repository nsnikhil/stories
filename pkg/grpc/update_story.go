package grpc

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"go.uber.org/zap"
)

func (ss *storiesServer) UpdateStory(ctx context.Context, req *proto.UpdateStoryRequest) (*proto.UpdateStoryResponse, error) {
	if err := validateRequestStory(ss.deps.cfg.GetBlogConfig(), req.GetStory(), true); err != nil {
		return &proto.UpdateStoryResponse{Success: false}, logAndGetError(err, "UpdateStory", "validateRequestStory", ss.deps.logger)
	}

	st, err := toDomainStory(req.GetStory())
	if err != nil {
		return &proto.UpdateStoryResponse{Success: false}, logAndGetError(err, "UpdateStory", "getStory", ss.deps.logger)
	}

	_, err = ss.deps.svc.GetStoriesService().UpdateStory(st)
	if err != nil {
		return &proto.UpdateStoryResponse{Success: false}, logAndGetError(err, "UpdateStory", "UpdateStory", ss.deps.logger)
	}

	ss.deps.logger.Info("story updated successfully", zap.String("method", "UpdateStory"))
	return &proto.UpdateStoryResponse{Success: true}, nil
}
