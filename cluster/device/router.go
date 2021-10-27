package device

import (
	"context"

	"github.com/aceaura/libra/magic"
)

type Router struct {
	name    string
	devices map[string]Device
	gateway Device
}

func (r *Router) String() string {
	return r.name
}

func (r *Router) LinkGateway(device Device) {
	r.gateway = device
}

func (r *Router) Process(ctx context.Context, route Route, data []byte) error {
	if route.Taking() {
		return r.gateway.Process(ctx, route, data)
	}
	return r.localProcess(ctx, route.Forward(), data)
}

func (r *Router) Link(device Device) {
	name := standardize(device.String(), magic.SeparatorNone)
	if name == "" {
		name = reflectTypeName(device)
	}
	r.devices[name] = device
}

func (r *Router) localProcess(ctx context.Context, route Route, data []byte) error {
	name := route.Name()
	device, ok := r.devices[name]
	if !ok {
		return route.Error(ErrRouteMissingDevice)
	}
	return device.Process(ctx, route, data)
}
