package server

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ss *storiesServer) GetMostViewedStories(context.Context, *proto.MostViewedStoriesRequest) (*proto.MostViewedStoriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMostViewedStories not implemented")
}
