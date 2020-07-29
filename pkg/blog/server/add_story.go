package server

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/pkg/blog/domain"
	"go.uber.org/zap"
)

func (ss *storiesServer) AddStory(ctx context.Context, req *proto.AddStoryRequest) (*proto.AddStoryResponse, error) {
	if err := validateRequestStory(ss.deps.cfg.GetBlogConfig(), req.GetStory(), false); err != nil {
		return &proto.AddStoryResponse{Success: false}, logAndGetError(err, "AddStory", "validateRequestStory", ss.deps.logger)
	}

	st, err := domain.NewVanillaStory(req.GetStory().GetTitle(), req.GetStory().GetBody())
	if err != nil {
		return &proto.AddStoryResponse{Success: false}, logAndGetError(err, "AddStory", "NewVanillaStory", ss.deps.logger)
	}

	if err := ss.deps.svc.GetStoriesService().AddStory(st); err != nil {
		return &proto.AddStoryResponse{Success: false}, logAndGetError(err, "AddStory", "AddStory", ss.deps.logger)
	}

	ss.deps.logger.Info("story added successfully", zap.String("method", "AddStory"))
	return &proto.AddStoryResponse{Success: true}, nil
}
