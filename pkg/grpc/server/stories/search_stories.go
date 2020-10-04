package stories

import (
	"context"
	"errors"
	"github.com/nsnikhil/stories-proto/proto"
)

func (ss *Server) SearchStories(ctx context.Context, req *proto.SearchStoriesRequest) (*proto.SearchStoriesResponse, error) {
	return &proto.SearchStoriesResponse{}, errors.New("UNIMPLEMENTED")
}
