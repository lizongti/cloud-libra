package service

import (
	"context"
	"reflect"

	"github.com/aceaura/libra/cluster"
)

var (
	typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()
	typeOfError   = reflect.TypeOf((*error)(nil)).Elem()
	typeOfBytes   = reflect.TypeOf(([]byte)(nil))
)

type Service interface {
	Bind(Router)
}

type service struct {
	codec      cluster.Codec
	renameFunc func(reflect.Method) string
	codeFunc   func(reflect.Method) uint64
}

type ServiceBase service

func (s *service) Bind(router *Router) {
	for _, handler := range s.extractHandlers() {
		handler.Register(router)
	}
}

func (s *service) OnInit() {

}

func (s *service) extractHandlers() []*Handler {
	t := reflect.TypeOf(s)
	handlers := make([]*Handler, 0)
	for index := 0; index < t.NumMethod(); index++ {
		method := t.Method(index)
		if !s.isMethodHandler(method) {
			continue
		}

		handler := NewHandler()

		mt := method.Type
		if mt.In(2) != typeOfBytes {
			handler.WithCodec(s.codec)
		}

		if s.renameFunc != nil {
			handler.WithName(s.renameFunc(method))
		}

		if s.codeFunc != nil {
			handler.WithCode(s.codeFunc(method))
		}
		handlers = append(handlers, 0)
	}
	return handlers
}

func (s *service) isMethodHandler(method reflect.Method) bool {
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
	if t := mt.In(1); !t.Implements(typeOfContext) {
		return false
	}

	// Check error
	if t := mt.Out(1); !t.Implements(typeOfBytes) {
		return false
	}

	// Check request:  pointer or bytes
	if t := mt.In(2); t.Kind() != reflect.Ptr && t != typeOfBytes {
		return false
	}

	// Check response: pointer or bytes
	if t := mt.Out(0); t.Kind() != reflect.Ptr && t != typeOfBytes {
		return false
	}

	return true
}
