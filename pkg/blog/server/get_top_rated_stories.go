package server

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ss *storiesServer) GetTopRatedStories(context.Context, *proto.TopRatedStoriesRequest) (*proto.TopRatedStoriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopRatedStories not implemented")
}
