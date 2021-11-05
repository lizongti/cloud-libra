package device

import (
	"context"
	"math/rand"
	"sync"

	"github.com/aceaura/libra/magic"
)

type Base struct {
	devices map[string][]Device
	gateway Device
	rwMutex sync.RWMutex
}

var _ Device = (*Base)(nil)
var empty Device = NewBase()

func NewBase() *Base {
	return &Base{
		devices: make(map[string][]Device),
	}
}

func (b *Base) String() string {
	return ""
}

func (b *Base) Tree() string {
	return ""
}

func (b *Base) Access(device Device) {
	b.gateway = device
}

func (b *Base) Gateway() Device {
	return b.gateway
}

func (b *Base) Extend(device Device) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	name := device.String()
	if name == "" {
		name = magic.TypeName(device)
	}

	for _, d := range b.devices[name] {
		if d == device {
			return
		}
	}
	b.devices[name] = append(b.devices[name], device)
}

func (b *Base) Route(name string) Device {
	b.rwMutex.RLock()
	defer b.rwMutex.RUnlock()

	devices, ok := b.devices[name]
	if !ok {
		return nil
	}
	return devices[rand.Intn(len(devices))]
}

func (b *Base) Process(_ context.Context, route Route, _ []byte) error {
	return route.Error(ErrRouteDeadEnd)
}
