package device

import (
	"context"
	"reflect"

	"github.com/aceaura/libra/cluster/component"
	"github.com/aceaura/libra/encoding"
	"github.com/aceaura/libra/magic"
	"github.com/aceaura/libra/scheduler"
)

type Service struct {
	component     component.Component
	encoding      encoding.Codec
	schedulerFunc func(context.Context) *scheduler.Scheduler
	handlers      map[string]*Handler
	gateway       Device
}

func NewService(opts ...serviceOpt) *Service {
	s := &Service{
		encoding: encoding.Emtpy(),
		handlers: make(map[string]*Handler),
		schedulerFunc: func(_ context.Context) *scheduler.Scheduler {
			return scheduler.Default()
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Service) String() string {
	return reflectTypeName(s.component)
}

func (s *Service) LinkGateway(device Device) {
	s.gateway = device
}

func (s *Service) Process(ctx context.Context, route Route, data []byte) error {
	if route.Taking() {
		return s.gateway.Process(ctx, route, data)
	}
	return s.localProcess(ctx, route.Forward(), data)
}

func (s *Service) localProcess(ctx context.Context, route Route, data []byte) error {
	name := route.Name()
	handler, ok := s.handlers[name]
	if !ok {
		return route.Error(ErrRouteMissingDevice)
	}
	return handler.Process(ctx, route, data)
}

func (s *Service) bind(component component.Component) {
	s.component = component
	s.handlers = make(map[string]*Handler)

	t := reflect.TypeOf(component)
	for index := 0; index < t.NumMethod(); index++ {
		method := t.Method(index)
		if !isMethodHandler(method) {
			continue
		}

		h := NewHandler().WithMethod(method)
		h.LinkGateway(s)
		s.handlers[h.String()] = h
	}
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

func reflectTypeName(i interface{}) string {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
		return reflect.TypeOf(i).Elem().Name()
	} else if v.Kind() == reflect.Struct {
		return reflect.TypeOf(i).Name()
	}
	return ""
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

func (serviceOption) WithEncoding(encoding encoding.Codec) serviceOpt {
	return func(s *Service) {
		s.WithEncoding(encoding)
	}
}

func (s *Service) WithEncoding(encoding encoding.Codec) *Service {
	s.encoding = encoding
	return s
}
