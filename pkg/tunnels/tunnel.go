package tunnels

type DeviceRequest struct {
	Pattern string
	Data string
}

type DeviceResponse struct {
	Data string
}

type DeviceTunnel interface {
	SendDeviceRequest(req DeviceRequest) (*DeviceResponse, error)
	Close() error
}
