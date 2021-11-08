package device

import (
	"context"
	"reflect"

	"github.com/aceaura/libra/cluster/component"
	"github.com/aceaura/libra/encoding"
	"github.com/aceaura/libra/magic"
	"github.com/aceaura/libra/scheduler"
)

type DispatchFunc func(context.Context, Route) *scheduler.Scheduler

var defaultDispatchFunc DispatchFunc

func init() {
	defaultDispatchFunc = func(context.Context, Route) *scheduler.Scheduler {
		return scheduler.Default()
	}
}

type Service struct {
	*Base
	component    component.Component
	encoding     encoding.Encoding
	dispatchFunc DispatchFunc
}

func NewService(opts ...serviceOpt) *Service {
	s := &Service{
		Base:         NewBase(),
		encoding:     encoding.Nil(),
		dispatchFunc: defaultDispatchFunc,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Service) String() string {
	return magic.TypeName(s.component)
}

func (s *Service) Process(ctx context.Context, route Route, data []byte) error {
	if route.Assembling() {
		return s.gateway.Process(ctx, route, data)
	}
	return s.localProcess(ctx, route.Forward(), data)
}

func (s *Service) localProcess(ctx context.Context, route Route, data []byte) error {
	name := route.Name()
	device := s.Route(name)
	if device == nil {
		return route.Error(ErrRouteMissingDevice)
	}
	return device.Process(ctx, route, data)
}

func (s *Service) bind(component component.Component) {
	s.component = component

	t := reflect.TypeOf(component)
	for index := 0; index < t.NumMethod(); index++ {
		method := t.Method(index)
		if !isMethodHandler(method) {
			continue
		}

		h := NewHandler().WithMethod(method)
		h.Access(s)
		s.Extend(h)
	}
}

func (s *Service) dispatch(ctx context.Context, route Route) *scheduler.Scheduler {
	return s.dispatchFunc(ctx, route)
}

func isMethodHandler(method reflect.Method) bool {
	mt := method.Type
	// Check method is exported
	if mt.PkgPath() != "" {
		return false
	}

	// Check num in
	if mt.NumIn() != 3 {
		return false
	}

	// Check num out
	if mt.NumOut() != 2 {
		return false
	}

	// Check context.Context
	if t := mt.In(1); !t.Implements(magic.TypeOfContext) {
		return false
	}

	// Check error
	if t := mt.Out(1); !t.Implements(magic.TypeOfError) {
		return false
	}

	// Check request:  pointer or bytes
	if t := mt.In(2); t.Kind() != reflect.Ptr && t != magic.TypeOfBytes {
		return false
	}

	// Check response: pointer or bytes
	if t := mt.Out(0); t.Kind() != reflect.Ptr && t != magic.TypeOfBytes {
		return false
	}

	return true
}

type serviceOpt func(*Service)
type serviceOption struct{}

var ServiceOption serviceOption

func (serviceOption) WithComponent(component component.Component) serviceOpt {
	return func(s *Service) {
		s.WithComponent(component)
	}
}

func (s *Service) WithComponent(component component.Component) *Service {
	s.bind(component)
	return s
}

func (serviceOption) WithEncoding(encoding encoding.Encoding) serviceOpt {
	return func(s *Service) {
		s.WithEncoding(encoding)
	}
}

func (s *Service) WithEncoding(encoding encoding.Encoding) *Service {
	s.encoding = encoding
	return s
}

func (s *Service) DispatchFunc() {

}
