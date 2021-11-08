package device

import (
	"context"
	"sync"

	"github.com/aceaura/libra/core/message"
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

func (r *Router) Process(ctx context.Context, msg *message.Message) error {
	if msg.State() == message.MessageStateAssembling {
		return r.gateway.Process(ctx, msg)
	}
	return r.localProcess(ctx, msg.Forward())
}

func (r *Router) localProcess(ctx context.Context, msg *message.Message) error {
	device := r.Locate(msg.Position())
	if device == nil {
		return routeErr(msg.RouteString(), ErrRouteMissingDevice)
	}
	return device.Process(ctx, msg)
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
