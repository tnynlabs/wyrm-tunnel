package tunnels

import (
	"sync"
	"errors"
)

type Registry interface {
	Add(deviceID int64, tun DeviceTunnel) error
	Get(deviceID int64) (DeviceTunnel, error)
	Delete(deviceID int64) error
}

type mapRegistry struct {
	sync.Mutex

	tunnelMap map[int64]DeviceTunnel
}

// CreateRegistry creates a new instance of map based tunnel registry
func CreateRegistry() Registry {
	return &mapRegistry{
		tunnelMap: make(map[int64]DeviceTunnel),
	}
}

func (reg *mapRegistry) Add(deviceID int64, tun DeviceTunnel) error {
	reg.Lock()
	defer reg.Unlock()

	reg.tunnelMap[deviceID] = tun

	return nil
}

func (reg *mapRegistry) Get(deviceID int64) (DeviceTunnel, error) {
	reg.Lock()
	defer reg.Unlock()

	tun, ok := reg.tunnelMap[deviceID]
	if !ok {
		return nil, errors.New("Not Found")
	}

	return tun, nil
}

func (reg *mapRegistry) Delete(deviceID int64) error {
	reg.Lock()
	defer reg.Unlock()

	if _, ok := reg.tunnelMap[deviceID]; !ok {
		return errors.New("Not Found")
	}

	delete(reg.tunnelMap, deviceID)

	return nil
}
