package device

import (
	"context"
	"reflect"

	"github.com/aceaura/libra/magic"
	"github.com/aceaura/libra/scheduler"
)

type Handler struct {
	method  reflect.Method
	gateway Device
}

func NewHandler(opts ...handlerOpt) *Handler {
	h := &Handler{}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *Handler) String() string {
	return h.method.Name
}

func (h *Handler) LinkGateway(device Device) {
	h.gateway = device
}

func (h *Handler) Process(ctx context.Context, route Route, data []byte) error {
	if route.Taking() {
		return h.gateway.Process(ctx, route, data)
	}
	return h.localProcess(ctx, route, data)
}

func (h *Handler) localProcess(ctx context.Context, route Route, reqData []byte) error {
	if h.method.Type == magic.TypeNil {
		return nil
	}

	s := h.gateway.(*Service)

	stage := func(t *scheduler.Task) error {
		respData, err := h.do(ctx, reqData)
		if err != nil {
			return err
		}
		return s.Process(ctx, route.Reverse(), respData)
	}
	scheduler.NewTask(
		scheduler.TaskOption.WithContext(ctx),
		scheduler.TaskOption.WithStage(stage),
	).Publish(s.schedulerFunc(ctx))
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

type handlerOpt func(*Handler)
type handlerOption struct{}

var HandlerOption handlerOption

func (handlerOption) WithMethod(method reflect.Method) handlerOpt {
	return func(h *Handler) {
		h.WithMethod(method)
	}
}

func (h *Handler) WithMethod(method reflect.Method) *Handler {
	h.method = method
	return h
}
