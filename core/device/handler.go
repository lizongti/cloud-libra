package device

import (
	"context"
	"reflect"

	"github.com/aceaura/libra/core/message"
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

func (h *Handler) Process(ctx context.Context, msg *message.Message) error {
	if msg.State() == message.MessageStateAssembling {
		return h.gateway.Process(ctx, msg)
	}
	return h.localProcess(ctx, msg)
}

func (h *Handler) localProcess(ctx context.Context, reqMsg *message.Message) error {
	if h.method.Type == magic.TypeNil {
		return nil
	}

	s := h.gateway.(*Service)
	scheduler.NewTask(
		scheduler.TaskOption.WithContext(ctx),
		scheduler.TaskOption.WithStage(func(t *scheduler.Task) error {
			ctx := t.Context()
			respMsg, err := h.do(ctx, reqMsg)
			if err != nil {
				return err
			}
			return s.Process(ctx, respMsg)
		}),
	).Publish(s.dispatch(ctx, reqMsg))

	return nil
}

func (h *Handler) do(ctx context.Context, reqMsg *message.Message) (*message.Message, error) {
	s := h.gateway.(*Service)
	mt := h.method.Type
	reqData := reqMsg.Data()
	req := reflect.New(mt.In(2).Elem()).Interface()
	err := s.encoding.Unmarshal(reqData, req)
	if err != nil {
		return nil, err
	}

	in := []reflect.Value{reflect.ValueOf(s.component), reflect.ValueOf(ctx), reflect.ValueOf(req)}

	out := h.method.Func.Call(in)
	if e := out[1].Interface(); e != nil {
		return nil, e.(error)
	}
	resp := out[0].Interface()
	respData, err := s.encoding.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return reqMsg.Reply(respData), nil
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
