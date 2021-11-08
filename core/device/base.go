package device

import (
	"context"
	"math/rand"
	"sync"

	"github.com/aceaura/libra/core/route"
	"github.com/aceaura/libra/magic"
)

type Base struct {
	extensions map[string][]Device
	gateway    Device
	rwMutex    sync.RWMutex
}

var _ Device = (*Base)(nil)
var empty Device = NewBase()

func NewBase() *Base {
	return &Base{
		extensions: make(map[string][]Device),
	}
}

func (b *Base) String() string {
	return magic.Anonymous
}

func (b *Base) Access(device Device) {
	b.gateway = device
}

func (b *Base) Gateway() Device {
	return b.gateway
}

func (b *Base) Extensions() map[string][]Device {
	return b.extensions
}

func (b *Base) Extend(device Device) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	name := device.String()

	for _, d := range b.extensions[name] {
		if d == device {
			return
		}
	}
	b.extensions[name] = append(b.extensions[name], device)
}

func (b *Base) Route(name string) Device {
	b.rwMutex.RLock()
	defer b.rwMutex.RUnlock()

	devices, ok := b.extensions[name]
	if !ok {
		return nil
	}
	return devices[rand.Intn(len(devices))]
}

func (b *Base) Process(_ context.Context, r route.Route, _ []byte) error {
	return r.Error(route.ErrRouteDeadEnd)
}
