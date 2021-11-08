package device

import (
	"context"
	"reflect"

	"github.com/aceaura/libra/core/route"
	"github.com/aceaura/libra/core/scheduler"
	"github.com/aceaura/libra/magic"
)

type Handler struct {
	*Base
	method reflect.Method
}

func NewHandler(opts ...funcHandlerOption) *Handler {
	h := &Handler{
		Base: NewBase(),
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *Handler) String() string {
	return h.method.Name
}

func (h *Handler) Process(ctx context.Context, rt route.Route, data []byte) error {
	if rt.Assembling() {
		return h.gateway.Process(ctx, rt, data)
	}
	return h.localProcess(ctx, rt, data)
}

func (h *Handler) localProcess(ctx context.Context, rt route.Route, reqData []byte) error {
	if h.method.Type == magic.TypeNil {
		return nil
	}

	s := h.gateway.(*Service)
	scheduler.NewTask(
		scheduler.TaskOption.WithContext(ctx),
		scheduler.TaskOption.WithStage(func(t *scheduler.Task) error {
			ctx := t.Context()
			respData, err := h.do(ctx, reqData)
			if err != nil {
				return err
			}
			return s.Process(ctx, rt.Reverse(), respData)
		}),
	).Publish(s.dispatch(ctx, rt))

	return nil
}

func (h *Handler) do(ctx context.Context, reqData []byte) (respData []byte, err error) {
	s := h.gateway.(*Service)
	mt := h.method.Type
	var req interface{}
	if mt.In(2) == magic.TypeOfBytes {
		req = reqData
	} else {
		req = reflect.New(mt.In(2).Elem()).Interface()
		if err = s.encoding.Unmarshal(reqData, req); err != nil {
			return nil, err
		}
	}

	in := []reflect.Value{reflect.ValueOf(s.component), reflect.ValueOf(ctx), reflect.ValueOf(req)}

	out := h.method.Func.Call(in)
	if e := out[1].Interface(); e != nil {
		return nil, e.(error)
	}
	resp := out[0].Interface()
	if respData, err = s.encoding.Marshal(resp); err != nil {
		return nil, err
	}
	return respData, err
}

type funcHandlerOption func(*Handler)
type handlerOption struct{}

var HandlerOption handlerOption

func (handlerOption) WithMethod(method reflect.Method) funcHandlerOption {
	return func(h *Handler) {
		h.WithMethod(method)
	}
}

func (h *Handler) WithMethod(method reflect.Method) *Handler {
	h.method = method
	return h
}
