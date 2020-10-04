package stories

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
)

func (ss *Server) GetStory(ctx context.Context, req *proto.GetStoryRequest) (*proto.GetStoryResponse, error) {
	st, err := ss.svc.GetStory(req.GetStoryID())
	if err != nil {
		return nil, err
	}

	return &proto.GetStoryResponse{Story: toProtoStory(st)}, nil
}
