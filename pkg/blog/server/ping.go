package server

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
)

func (ss *storiesServer) Ping(ctx context.Context, in *proto.PingRequest) (*proto.PingResponse, error) {
	ss.logger.Info("[storiesServer] [Ping]")
	return &proto.PingResponse{Message: "pong"}, nil
}




