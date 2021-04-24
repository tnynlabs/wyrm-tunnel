package manager

import (
	"net"
	"context"
	
	pb "github.com/tnynlabs/wyrm-tunnel/pkg/manager/protobuf"
	"github.com/tnynlabs/wyrm-tunnel/pkg/tunnels"
	
	"google.golang.org/grpc"
	"github.com/golang/protobuf/ptypes/empty"
)

type tunnelManagerServer struct {
	pb.UnimplementedTunnelManagerServer

	registry tunnels.Registry
}

func NewServer(registry tunnels.Registry) pb.TunnelManagerServer {
	return &tunnelManagerServer{
		registry: registry,
	}
}

func RunServer(address string, server pb.TunnelManagerServer) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	// TODO: Check needed server options
	grpcServer := grpc.NewServer()
	pb.RegisterTunnelManagerServer(grpcServer, server)

	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func (s *tunnelManagerServer) RevokeDevice(ctx context.Context, req *pb.RevokeRequest) (*empty.Empty, error) {
	deviceID := int64(req.GetDeviceId())
	deviceTunnel, err := s.registry.Get(deviceID)
	if err != nil {
		return nil, err
	}

	err = deviceTunnel.Close()
	if err != nil {
		return nil, err
	}

	err = s.registry.Delete(deviceID)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *tunnelManagerServer) InvokeDevice(ctx context.Context, req *pb.InvokeRequest) (*pb.InvokeResponse, error) {
	deviceID := int64(req.GetDeviceId())
	deviceTunnel, err := s.registry.Get(deviceID)
	if err != nil {
		return nil, err
	}

	tunnelReq := tunnels.DeviceRequest{
		Pattern: req.GetPattern(),
		Data: req.GetData(),
	}

	tunnelResp, err := deviceTunnel.SendDeviceRequest(tunnelReq)
	if err != nil {
		return nil, err
	}

	resp := pb.InvokeResponse{
		Data: tunnelResp.Data,
	}

	return &resp, nil
}
