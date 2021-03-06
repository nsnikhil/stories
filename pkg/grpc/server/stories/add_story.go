package stories

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/pkg/liberr"
	"github.com/nsnikhil/stories/pkg/story/model"
)

func (ss *Server) AddStory(ctx context.Context, req *proto.AddStoryRequest) (*proto.AddStoryResponse, error) {
	st, err := model.NewStoryBuilder().
		SetTitle(ss.cfg.TitleMaxLength(), req.GetStory().GetTitle()).
		SetBody(ss.cfg.BodyMaxLength(), req.GetStory().GetBody()).
		Build()

	if err != nil {
		return &proto.AddStoryResponse{Success: false},
			liberr.WithArgs(liberr.Operation("Server.AddStory"), err)
	}

	if err := ss.svc.AddStory(st); err != nil {
		return &proto.AddStoryResponse{Success: false},
			liberr.WithArgs(liberr.Operation("Server.AddStory"), err)
	}

	//TODO: ADD SUCCESS LOG
	return &proto.AddStoryResponse{Success: true}, nil
}
