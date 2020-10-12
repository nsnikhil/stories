package stories

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
)

func (ss *Server) Ping(ctx context.Context, in *proto.PingRequest) (*proto.PingResponse, error) {
	//TODO: ADD SUCCESS LOG
	return &proto.PingResponse{Message: "pong"}, nil
}
