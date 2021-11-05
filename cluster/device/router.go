package device

import (
	"context"
	"sync"
)

type Router struct {
	*Base
	name    string
	rwMutex sync.RWMutex
}

func NewRouter(opts ...routerOpt) *Router {
	r := &Router{
		Base: NewBase(),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Router) String() string {
	return r.name
}

func (r *Router) Process(ctx context.Context, route Route, data []byte) error {
	if route.Assembling() {
		return r.gateway.Process(ctx, route, data)
	}
	return r.localProcess(ctx, route.Forward(), data)
}

func (r *Router) localProcess(ctx context.Context, route Route, data []byte) error {
	name := route.Name()
	device := r.Route(name)
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
		r.Extend(device)
		device.Access(r)
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
