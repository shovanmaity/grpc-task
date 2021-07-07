package server

import (
	"context"

	generated "github.com/shovanmaity/grpc-task/gen/go"
)

var _ generated.PingServiceServer = &PingServer{}

type PingServer struct {
	generated.UnimplementedPingServiceServer
}

func NewPingServer() *PingServer {
	return &PingServer{}
}

func (s *PingServer) Ping(context.Context, *generated.EmptyMessage) (*generated.EmptyMessage, error) {
	return &generated.EmptyMessage{}, nil
}
