package stories

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/pkg/liberr"
)

func (ss *Server) UpdateStory(ctx context.Context, req *proto.UpdateStoryRequest) (*proto.UpdateStoryResponse, error) {
	st, err := toDomainStory(ss.cfg, req.GetStory())
	if err != nil {
		return &proto.UpdateStoryResponse{Success: false}, liberr.WithArgs(liberr.Operation("Server.UpdateStory"), err)
	}

	_, err = ss.svc.UpdateStory(st)
	if err != nil {
		return &proto.UpdateStoryResponse{Success: false}, liberr.WithArgs(liberr.Operation("Server.UpdateStory"), err)
	}

	//TODO: ADD SUCCESS LOG
	return &proto.UpdateStoryResponse{Success: true}, nil
}
