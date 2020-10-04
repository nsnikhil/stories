package stories

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
)

func (ss *Server) DeleteStory(ctx context.Context, req *proto.DeleteStoryRequest) (*proto.DeleteStoryResponse, error) {
	_, err := ss.svc.DeleteStory(req.GetStoryID())
	if err != nil {
		return &proto.DeleteStoryResponse{Success: false}, err
	}

	return &proto.DeleteStoryResponse{Success: true}, nil
}
