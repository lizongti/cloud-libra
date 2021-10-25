package device

import (
	"context"
	"reflect"
)

type Router struct {
	Device
	devices map[string]Device
	gateway Device
}

func (r *Router) String() string {
	return reflect.Indirect(reflect.ValueOf(r)).Type().Name()
}

func (r *Router) Gateway(device Device) {
	r.gateway = device
}

func (r *Router) Process(ctx context.Context, route Route, data []byte) error {
	deviceType := route.deviceType()
	if deviceType == DeviceTypeBus {
		return r.gateway.Process(ctx, route, data)
	} else if deviceType == DeviceTypeRouter {
		return r.localProcess(ctx, route, data)
	}
	return ErrRouteDeadEnd
}

func (r *Router) Link(device Device) {
	// TODO
	r.devices[device.String()] = device
}

func (r *Router) localProcess(ctx context.Context, route Route, data []byte) error {
	name := route.deviceName()
	device, ok := r.devices[name]
	if !ok {
		return ErrRouteMissingDevice
	}
	return device.Process(ctx, route.forward(), data)
}
