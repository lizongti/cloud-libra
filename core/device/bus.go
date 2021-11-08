package device

import (
	"context"
	"sync"

	"github.com/aceaura/libra/core/route"
	"github.com/aceaura/libra/magic"
)

type Bus struct {
	*Base
	rwMutex sync.RWMutex
}

func NewBus(opts ...funcBusOption) *Bus {
	b := &Bus{
		Base: NewBase(),
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *Bus) String() string {
	return magic.TypeName(b)
}

func (b *Bus) Access(device Device) {
	b.gateway = Hole()
}

func (b *Bus) Process(ctx context.Context, rt route.Route, data []byte) error {
	if rt.Assembling() {
		return b.localProcess(ctx, rt.Forward(), data)
	}

	return rt.Error(route.ErrRouteDeadEnd)
}

func (b *Bus) localProcess(ctx context.Context, rt route.Route, data []byte) error {
	name := rt.Name()
	device := b.Route(name)
	if device == nil {
		return rt.Error(route.ErrRouteMissingDevice)
	}
	return device.Process(ctx, rt, data)
}

type funcBusOption func(*Bus)
type busOption struct{}

var BusOption busOption

func (busOption) WithDevice(devices ...Device) funcBusOption {
	return func(b *Bus) {
		b.WithDevice(devices...)
	}
}

func (b *Bus) WithDevice(devices ...Device) *Bus {
	for _, device := range devices {
		b.Extend(device)
		device.Access(b)
	}
	return b
}
