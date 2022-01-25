package device

import (
	"context"

	"github.com/aceaura/libra/core/message"
)

type Router struct {
	*Base
	opts routerOptions
}

func NewRouter(opt ...funcRouterOption) *Router {
	opts := defaultRouterOptions

	for _, o := range opt {
		o.apply(&opts)
	}

	r := &Router{
		Base: NewBase(),
	}

	r.link()

	return r
}

var bus *Router = NewRouter().WithBus().WithName("Bus")

func Bus() *Router {
	return bus
}

func (r *Router) String() string {
	return r.opts.name
}

func (r *Router) Process(ctx context.Context, msg *message.Message) error {
	if r.opts.bus {
		if !msg.Route.Dispatching() {
			msg.Route = msg.Route.Forward()
			return r.localProcess(ctx, msg)
		}

		return msg.Route.Error(ErrRouteDeadEnd)
	}

	if !msg.Route.Dispatching() {
		if r.gateway == nil {
			return ErrGatewayNotFound
		}
		return r.gateway.Process(ctx, msg)
	}
	msg.Route = msg.Route.Forward()
	return r.localProcess(ctx, msg)
}

func (r *Router) localProcess(ctx context.Context, msg *message.Message) error {
	device := r.Locate(msg.Route.Position())
	if device == nil {
		return msg.Route.Error(ErrRouteMissingDevice)
	}
	return device.Process(ctx, msg)
}

func (r *Router) link() {
	for _, d := range r.opts.devices {
		d.Access(r)
		r.Extend(d)
	}
}

type routerOptions struct {
	name    string
	bus     bool
	devices []Device
}

var defaultRouterOptions = routerOptions{
	name: "",
	bus:  false,
}

type ApplyRouterOption interface {
	apply(*routerOptions)
}

type funcRouterOption func(*routerOptions)

func (f funcRouterOption) apply(opt *routerOptions) {
	f(opt)
}

type routerOption int

var RouterOption routerOption

func (routerOption) Device(devices ...Device) funcRouterOption {
	return func(r *routerOptions) {
		r.devices = append(r.devices, devices...)
	}
}

func (r *Router) WithDevice(devices ...Device) *Router {
	RouterOption.Device(devices...).apply(&r.opts)
	r.link()
	return r
}

func (routerOption) Service(services ...interface{}) funcRouterOption {
	return func(r *routerOptions) {
		for _, service := range services {
			r.devices = append(r.devices, extractHandlers(service)...)
		}
	}
}

func (r *Router) WithService(services ...interface{}) *Router {
	RouterOption.Service(services...).apply(&r.opts)
	r.link()
	return r
}

func (routerOption) Name(name string) funcRouterOption {
	return func(r *routerOptions) {
		r.name = name
	}
}

func (r *Router) WithName(name string) *Router {
	RouterOption.Name(name).apply(&r.opts)
	return r
}

func (routerOption) Bus() funcRouterOption {
	return func(r *routerOptions) {
		r.bus = true
	}
}

func (r *Router) WithBus() *Router {
	RouterOption.Bus().apply(&r.opts)
	return r
}
