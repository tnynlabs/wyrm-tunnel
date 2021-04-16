package manager

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

type Server struct {
	UnimplementedTunnelManagerServer
}

func (s *Server) RevokeDevice(ctx context.Context, info *DeviceInfo) (*empty.Empty, error) {
	return nil, nil
}

func (s *Server) InvokeDevice(ctx context.Context, r *InvokeRequest) (*InvokeResponse, error) {
	return nil, nil
}
