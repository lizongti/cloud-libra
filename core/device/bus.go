package device

import (
	"context"
	"sync"

	"github.com/aceaura/libra/core/message"
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

func (b *Bus) Process(ctx context.Context, msg *message.Message) error {
	if msg.State() == message.MessageStateAssembling {
		return b.localProcess(ctx, msg.Forward())
	}

	return routeErr(msg.RouteString(), ErrRouteDeadEnd)
}

func (b *Bus) localProcess(ctx context.Context, msg *message.Message) error {
	device := b.Locate(msg.Position())
	if device != nil {
		return device.Process(ctx, msg)
	}

	return routeErr(msg.RouteString(), ErrRouteMissingDevice)
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
