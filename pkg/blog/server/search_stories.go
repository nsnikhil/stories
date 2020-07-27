package server

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ss *storiesServer) SearchStories(context.Context, *proto.SearchStoriesRequest) (*proto.SearchStoriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchStories not implemented")
}
