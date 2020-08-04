package server

import (
	"context"
	"errors"
	"github.com/nsnikhil/stories-proto/proto"
	"go.uber.org/zap"
)

func (ss *storiesServer) DeleteStory(ctx context.Context, req *proto.DeleteStoryRequest) (*proto.DeleteStoryResponse, error) {
	c, err := ss.deps.svc.GetStoriesService().DeleteStory(req.GetStoryID())
	if err != nil {
		return &proto.DeleteStoryResponse{Success: false}, logAndGetError(err, "DeleteStory", "DeleteStory", ss.deps.logger)
	}

	if c <= 0 {
		err := errors.New("failed to delete story")
		return &proto.DeleteStoryResponse{Success: false}, logAndGetError(err, "DeleteStory", "DeleteStory", ss.deps.logger)
	}

	ss.deps.logger.Info("story deleted successfully", zap.String("method", "UpdateStory"))
	return &proto.DeleteStoryResponse{Success: true}, nil
}
