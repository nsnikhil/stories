package health

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
)

type Server struct {
	proto.UnimplementedHealthServer
}

func NewHealthServer() *Server {
	return &Server{}
}

func (hs *Server) Check(context.Context, *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{
		Status: proto.HealthCheckResponse_SERVING,
	}, nil
}

func (hs *Server) Watch(context.Context, *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{
		Status: proto.HealthCheckResponse_SERVING,
	}, nil
}
