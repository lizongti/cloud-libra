package main

import (
	"reflect"
	"strings"
	"sync"

	"github.com/lizongti/libra/module/core/component"
	"github.com/lizongti/libra/module/core/context"
	"github.com/lizongti/libra/module/core/encoding"
	"github.com/lizongti/libra/module/core/handler"
	"github.com/lizongti/libra/module/core/router"
)

type Router struct {
	component.ComponentBase
	m        sync.Map
	encoding encoding.Encoding
}

var _ router.Router = (*Router)(nil)

func (r *Router) Handle(route string, handler handler.Handler) {
	r.m.Store(route, handler)
}

func (r *Router) Serve(ctx context.Context, payload []byte) ([]byte, error) {
	route := strings.Split(strings.TrimLeft(ctx.GetLeftAddr(), "/"), "/")[0]
	v, ok := r.m.Load(route)
	if !ok {
		panic("route not found")
	}
	return v.(handler.Handler).Serve(ctx, payload)
}

func (r *Router) WithEncoding(encoding encoding.Encoding) {
	r.encoding = encoding
}

func (r *Router) OnInit() {
	r.installHandlers()
}

func (r *Router) installHandlers() {
	t := reflect.TypeOf(r)
	for m := 0; m < t.NumMethod(); m++ {
		method := t.Method(m)
		if method.Type.NumIn() == 2 && method.Type.NumOut() == 2 {
			handler := &Handler{
				method:   method,
				receiver: reflect.ValueOf(r),
				encoding: r.encoding,
			}
			r.Handle(handler.String(), handler)
		}
	}
}

func (r *Router) Extend(route string, rFront *Router) {
	r.Handle(route, rFront)
}
