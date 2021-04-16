package tunnels

type DeviceRequest struct {

}

type DeviceResponse struct {
	
}

type DeviceTunnel interface {
	SendDeviceRequest(r *DeviceRequest) (*DeviceResponse, error)
	Close() error // to be called when revoked
}
