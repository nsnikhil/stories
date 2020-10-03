package stories

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
)

func (ss *StoriesServer) Ping(ctx context.Context, in *proto.PingRequest) (*proto.PingResponse, error) {
	return &proto.PingResponse{Message: "pong"}, nil
}
