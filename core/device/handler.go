package device

import (
	"context"
	"reflect"

	"github.com/aceaura/libra/core/encoding"
	"github.com/aceaura/libra/core/message"
	"github.com/aceaura/libra/magic"
)

type Handler struct {
	*Base
	receiver reflect.Value
	method   reflect.Method
}

func (h *Handler) String() string {
	return h.method.Name
}

func (h *Handler) Process(ctx context.Context, msg *message.Message) error {
	if msg.Route.Assembling() {
		return h.gateway.Process(ctx, msg)
	}
	return h.localProcess(ctx, msg)
}

func (h *Handler) localProcess(ctx context.Context, reqMsg *message.Message) error {
	if h.method.Type == magic.TypeNil {
		return nil
	}

	respMsg, err := h.do(ctx, reqMsg)
	if err != nil {
		return err
	}
	if respMsg != nil {
		return h.Process(ctx, respMsg)
	}

	return nil
}

func (h *Handler) do(ctx context.Context, reqMsg *message.Message) (*message.Message, error) {
	mt := h.method.Type
	var req interface{}
	if mt.In(2) == magic.TypeOfBytes {
		bytes := &encoding.Bytes{}
		err := reqMsg.Encoding.Unmarshal(reqMsg.Data, bytes)
		if err != nil {
			return nil, err
		}
		req = bytes.Data
	} else {
		req = reflect.New(mt.In(2).Elem()).Interface()
		err := reqMsg.Encoding.Unmarshal(reqMsg.Data, req)
		if err != nil {
			return nil, err
		}
	}

	in := []reflect.Value{h.receiver, reflect.ValueOf(ctx), reflect.ValueOf(req)}

	out := h.method.Func.Call(in)
	if e := out[1].Interface(); e != nil {
		return nil, e.(error)
	}

	resp := out[0].Interface()
	if resp == nil {
		return nil, nil
	}

	respMsg := &message.Message{
		ID:       reqMsg.ID,
		Route:    reqMsg.Route.Reverse(),
		Encoding: reqMsg.Encoding.Reverse(),
	}
	respData, err := reqMsg.Encoding.Marshal(resp)
	if err != nil {
		return nil, err
	}
	respMsg.Data = respData
	return respMsg, nil
}
