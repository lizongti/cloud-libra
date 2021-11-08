package device

import (
	"context"
	"sync"

	"github.com/aceaura/libra/core/route"
)

type Router struct {
	*Base
	name    string
	rwMutex sync.RWMutex
}

func NewRouter(opts ...funcRouterOption) *Router {
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

func (r *Router) Process(ctx context.Context, rt route.Route, data []byte) error {
	if rt.Assembling() {
		return r.gateway.Process(ctx, rt, data)
	}
	return r.localProcess(ctx, rt.Forward(), data)
}

func (r *Router) localProcess(ctx context.Context, rt route.Route, data []byte) error {
	name := rt.Name()
	device := r.Route(name)
	if device == nil {
		return rt.Error(route.ErrRouteMissingDevice)
	}
	return device.Process(ctx, rt, data)
}

type funcRouterOption func(*Router)
type routerOption struct{}

var RouterOption routerOption

func (routerOption) WithDevice(devices ...Device) funcRouterOption {
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

func (routerOption) WithName(name string) funcRouterOption {
	return func(r *Router) {
		r.WithName(name)
	}
}

func (r *Router) WithName(name string) *Router {
	r.name = name
	return r
}
