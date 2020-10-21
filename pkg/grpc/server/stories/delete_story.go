package stories

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/pkg/liberr"
)

func (ss *Server) DeleteStory(ctx context.Context, req *proto.DeleteStoryRequest) (*proto.DeleteStoryResponse, error) {
	_, err := ss.svc.DeleteStory(req.GetStoryID())
	if err != nil {
		return &proto.DeleteStoryResponse{Success: false}, liberr.WithArgs(liberr.Operation("Server.DeleteStory"), err)
	}

	//TODO: ADD SUCCESS LOG
	return &proto.DeleteStoryResponse{Success: true}, nil
}
