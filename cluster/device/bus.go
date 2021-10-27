package device

import (
	"context"
	"math/rand"
	"sync"
)

type Bus struct {
	devices map[string][]Device
	gateway Device
	rwMutex sync.RWMutex
}

func NewBus(opts ...busOpt) *Bus {
	b := &Bus{
		devices: make(map[string][]Device),
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *Bus) String() string {
	return "Bus"
}

func (b *Bus) LinkGateway(device Device) {
	b.gateway = Empty()
}

func (b *Bus) Process(ctx context.Context, route Route, data []byte) error {
	deviceType := route.deviceType()
	if deviceType == DeviceTypeBus {
		return b.localProcess(ctx, route.forward(), data)
	}

	return ErrRouteDeadEnd
}

func (b *Bus) localProcess(ctx context.Context, route Route, data []byte) error {
	name := route.deviceName()
	device := b.mutexFindDevice(name)
	if device == nil {
		return ErrRouteMissingDevice
	}
	return device.Process(ctx, route, data)
}

func (b *Bus) mutexLinkDevice(device Device) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	name := device.String()
	for _, d := range b.devices[name] {
		if d == device {
			return
		}
	}
	b.devices[name] = append(b.devices[name], device)
}

func (b *Bus) mutexFindDevice(name string) Device {
	b.rwMutex.RLock()
	defer b.rwMutex.RUnlock()

	devices, ok := b.devices[name]
	if !ok {
		return nil
	}
	return devices[rand.Intn(len(devices))]
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
		b.mutexLinkDevice(device)
		device.LinkGateway(b)
	}
	return b
}
