package tungrpc

import (
	"log"
	"errors"
	"net"

	"github.com/tnynlabs/wyrm/pkg/devices"
	"github.com/tnynlabs/wyrm-tunnel/pkg/tunnels"
	pb "github.com/tnynlabs/wyrm-tunnel/pkg/transport/tungrpc/protobuf"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc"
)

type grpcDeviceTunnelServer struct {
	pb.UnimplementedDeviceTunnelServer

	registry tunnels.Registry
	deviceService devices.Service
}

func NewServer(registry tunnels.Registry, deviceService devices.Service) pb.DeviceTunnelServer {
	return &grpcDeviceTunnelServer{
		registry: registry,
		deviceService: deviceService,
	}
}

func RunServer(address string, server pb.DeviceTunnelServer) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	// TODO: Check needed server options
	grpcServer := grpc.NewServer()
	pb.RegisterDeviceTunnelServer(grpcServer, server)

	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func (s *grpcDeviceTunnelServer) CreateTunnel(stream pb.DeviceTunnel_CreateTunnelServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return errors.New("Cannot load metadata")
	}

	keyEntry := md.Get("authKey")
	if len(keyEntry) != 1 {
		return errors.New("'authKey' entry not found in metadata")
	}

	device, err := s.deviceService.GetByKey(keyEntry[0])
	if err != nil {
		return errors.New("Invalid device auth key")
	}

	tun := grpcDeviceTunnel{
		stream: stream,
		closeChan: make(chan bool),
	}

	log.Println(device.ID)
	err = s.registry.Add(device.ID, &tun)
	if err != nil {
		return errors.New("Unable to establish tunnel connection")
	}

	<-tun.closeChan // wait until close is called on tunnel

	return nil
}

type grpcDeviceTunnel struct {
	stream pb.DeviceTunnel_CreateTunnelServer
	closeChan chan bool
}

func (tun *grpcDeviceTunnel) SendDeviceRequest(req tunnels.DeviceRequest) (*tunnels.DeviceResponse, error) {
	pbReq := pb.DeviceRequest{
		Pattern: req.Pattern,
		Data: req.Data,
	}

	err := tun.stream.Send(&pbReq)
	if err != nil {
		return nil, err
	}

	// TODO: Drop responses not related to current request
	pbResp, err := tun.stream.Recv()
	if err != nil {
		return nil, err
	}

	resp := tunnels.DeviceResponse{
		Data: pbResp.Data,
	}

	return &resp, nil
}

func (tun *grpcDeviceTunnel) Close() error {
	tun.closeChan <- true
	return nil
}
