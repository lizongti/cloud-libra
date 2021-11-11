package device

import (
	"context"
	"reflect"

	"github.com/aceaura/libra/core/component"
	"github.com/aceaura/libra/core/message"
	"github.com/aceaura/libra/magic"
)

type Router struct {
	*Base
	name string
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
	if msg.Route.Assembling() {
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

func (r *Router) extract(c component.Component) {
	if r.name == "" {
		r.name = magic.TypeName(c)
	}
	t := reflect.TypeOf(c)

	for index := 0; index < t.NumMethod(); index++ {
		method := t.Method(index)
		mt := method.Type

		switch {
		case mt.PkgPath() != "": // Check method is exported
			continue
		case mt.NumIn() != 3: // Check num in
			continue
		case mt.NumOut() != 2: // Check num in
			continue
		case !mt.In(1).Implements(magic.TypeOfContext): // Check context.Context
			continue
		case !mt.Out(1).Implements(magic.TypeOfError): // Check error
			continue
		case mt.In(2).Kind() != reflect.Ptr && mt.In(2) != magic.TypeOfBytes: // Check request:  pointer or bytes
			continue
		case mt.Out(0).Kind() != reflect.Ptr && mt.Out(0) != magic.TypeOfBytes: // Check response: pointer or bytes
			continue
		}

		receiver := reflect.ValueOf(c)
		h := &Handler{
			Base:     NewBase(),
			receiver: receiver,
			method:   method,
		}
		h.Access(r)
		r.Extend(h)
	}
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

func (routerOption) WithService(components ...component.Component) funcRouterOption {
	return func(r *Router) {
		r.WithService(components...)
	}
}

func (r *Router) WithService(components ...component.Component) {
	for _, component := range components {
		r.extract(component)
	}
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
