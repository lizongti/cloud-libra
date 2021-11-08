package device

import (
	"context"
	"reflect"

	"github.com/aceaura/libra/core/component"
	"github.com/aceaura/libra/core/encoding"
	"github.com/aceaura/libra/core/route"
	"github.com/aceaura/libra/core/scheduler"
	"github.com/aceaura/libra/magic"
)

type Service struct {
	*Base
	component  component.Component
	encoding   encoding.Encoding
	dispatcher scheduler.Dispatcher
}

func NewService(opts ...funcServiceOption) *Service {
	s := &Service{
		Base:       NewBase(),
		encoding:   encoding.Empty(),
		dispatcher: scheduler.DefaultDispatcher(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Service) String() string {
	return magic.TypeName(s.component)
}

func (s *Service) Process(ctx context.Context, rt route.Route, data []byte) error {
	if rt.Assembling() {
		return s.gateway.Process(ctx, rt, data)
	}
	return s.localProcess(ctx, rt.Forward(), data)
}

func (s *Service) localProcess(ctx context.Context, rt route.Route, data []byte) error {
	name := rt.Name()
	device := s.Route(name)
	if device == nil {
		return rt.Error(route.ErrRouteMissingDevice)
	}
	return device.Process(ctx, rt, data)
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

func (s *Service) dispatch(ctx context.Context, rt route.Route) *scheduler.Scheduler {
	return s.dispatcher.Dispatch(ctx, rt)
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

type funcServiceOption func(*Service)
type serviceOption struct{}

var ServiceOption serviceOption

func (serviceOption) WithComponent(component component.Component) funcServiceOption {
	return func(s *Service) {
		s.WithComponent(component)
	}
}

func (s *Service) WithComponent(component component.Component) *Service {
	s.bind(component)
	return s
}

func (serviceOption) WithEncoding(encoding encoding.Encoding) funcServiceOption {
	return func(s *Service) {
		s.WithEncoding(encoding)
	}
}

func (s *Service) WithEncoding(encoding encoding.Encoding) *Service {
	s.encoding = encoding
	return s
}

func (serviceOption) WithDispatcher(dispatcher scheduler.Dispatcher) funcServiceOption {
	return func(s *Service) {
		s.WithDispatcher(dispatcher)
	}
}

func (s *Service) WithDispatcher(dispatcher scheduler.Dispatcher) *Service {
	s.dispatcher = dispatcher
	return s
}
