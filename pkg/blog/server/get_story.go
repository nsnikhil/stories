package server

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ss *storiesServer) GetStory(context.Context, *proto.GetStoryRequest) (*proto.GetStoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStory not implemented")
}
