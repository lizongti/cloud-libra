package device

import (
	"context"
	"math/rand"
	"sync"

	"github.com/aceaura/libra/magic"
)

type Router struct {
	name    string
	devices map[string][]Device
	gateway Device
	rwMutex sync.RWMutex
}

func NewRouter(opts ...routerOpt) *Router {
	r := &Router{
		devices: make(map[string][]Device),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
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

func (r *Router) mutexLinkDevice(device Device) {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	name := standardize(device.String(), magic.SeparatorNone)
	if name == "" {
		name = reflectTypeName(device)
	}
	for _, d := range r.devices[name] {
		if d == device {
			return
		}
	}
	r.devices[name] = append(r.devices[name], device)
}

func (r *Router) mutexFindDevice(name string) Device {
	r.rwMutex.RLock()
	defer r.rwMutex.RUnlock()

	devices, ok := r.devices[name]
	if !ok {
		return nil
	}
	return devices[rand.Intn(len(devices))]
}

func (r *Router) localProcess(ctx context.Context, route Route, data []byte) error {
	name := route.Name()
	device := r.mutexFindDevice(name)
	if device == nil {
		return route.Error(ErrRouteMissingDevice)
	}
	return device.Process(ctx, route, data)
}

type routerOpt func(*Router)
type routerOption struct{}

var RouterOption routerOption

func (routerOption) WithDevice(devices ...Device) routerOpt {
	return func(r *Router) {
		r.WithDevice(devices...)
	}
}

func (r *Router) WithDevice(devices ...Device) *Router {
	for _, device := range devices {
		r.mutexLinkDevice(device)
		device.LinkGateway(r)
	}
	return r
}

func (routerOption) WithName(name string) routerOpt {
	return func(r *Router) {
		r.WithName(name)
	}
}

func (r *Router) WithName(name string) *Router {
	r.name = name
	return r
}
