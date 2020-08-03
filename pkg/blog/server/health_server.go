package server

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
)

type healthServer struct {
	proto.UnimplementedHealthServer
}

func newHealthServer() *healthServer {
	return &healthServer{}
}

func (hs *healthServer) Check(context.Context, *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{
		Status: proto.HealthCheckResponse_SERVING,
	}, nil
}

func (hs *healthServer) Watch(context.Context, *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{
		Status: proto.HealthCheckResponse_SERVING,
	}, nil
}
