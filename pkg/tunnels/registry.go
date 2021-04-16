package tunnels

type Registry interface {
	Add(deviceID int, tun DeviceTunnel) error
	Get(deviceID int) (DeviceTunnel, error)
	Delete(deviceID int) error
}

type mapRegistry struct {
	tunnelMap map[int]DeviceTunnel
}

// CreateRegistry creates a new instance of map based tunnel registry
func CreateRegistry() Registry {
	return &mapRegistry{
		tunnelMap: make(map[int]DeviceTunnel),
	}
}

func (reg *mapRegistry) Add(deviceID int, tun DeviceTunnel) error {
	return nil
}

func (reg *mapRegistry) Get(deviceID int) (DeviceTunnel, error) {
	return nil, nil
}

func (reg *mapRegistry) Delete(deviceID int) error {
	return nil
}