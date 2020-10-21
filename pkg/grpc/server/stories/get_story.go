package stories

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/pkg/liberr"
)

func (ss *Server) GetStory(ctx context.Context, req *proto.GetStoryRequest) (*proto.GetStoryResponse, error) {
	st, err := ss.svc.GetStory(req.GetStoryID())
	if err != nil {
		return nil, liberr.WithArgs(liberr.Operation("Server.GetStory"), err)
	}

	//TODO: ADD SUCCESS LOG
	return &proto.GetStoryResponse{Story: toProtoStory(st)}, nil
}
