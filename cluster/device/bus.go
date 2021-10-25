package device

import "context"

type Bus struct {
	Device
	devices map[string]Device
	gateway Device
}

func (b *Bus) String() string {
	return ""
}

func (b *Bus) Gateway(device Device) {
	b.gateway = Empty()
}

func (b *Bus) Process(ctx context.Context, route Route, data []byte) error {
	deviceType := route.deviceType()
	if deviceType == DeviceTypeBus {
		return b.localProcess(ctx, route, data)
	}

	return ErrRouteDeadEnd
}

func (b *Bus) Discover() {
	// TODO
}

func (b *Bus) localProcess(ctx context.Context, route Route, data []byte) error {
	name := route.deviceName()
	device, ok := b.devices[name]
	if !ok {
		return ErrRouteMissingDevice
	}
	return device.Process(ctx, route, data)
}
