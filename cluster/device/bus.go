package device

import (
	"context"
	"sync"

	"github.com/aceaura/libra/magic"
)

type Bus struct {
	*Base
	rwMutex sync.RWMutex
}

func NewBus(opts ...busOpt) *Bus {
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

func (b *Bus) Process(ctx context.Context, route Route, data []byte) error {
	if route.Assembling() {
		return b.localProcess(ctx, route.Forward(), data)
	}

	return route.Error(ErrRouteDeadEnd)
}

func (b *Bus) localProcess(ctx context.Context, route Route, data []byte) error {
	name := route.Name()
	device := b.Route(name)
	if device == nil {
		return route.Error(ErrRouteMissingDevice)
	}
	return device.Process(ctx, route, data)
}

type busOpt func(*Bus)
type busOption struct{}

var BusOption busOption

func (busOption) WithDevice(devices ...Device) busOpt {
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
