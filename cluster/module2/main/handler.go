package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lizongti/libra/module/core/context"
	"github.com/lizongti/libra/module/core/encoding"
	"github.com/lizongti/libra/module/core/handler"
)

type Handler struct {
	method   reflect.Method
	receiver reflect.Value
	encoding encoding.Encoding
}

var _ handler.Handler = (*Handler)(nil)

func (h *Handler) Serve(ctx context.Context, inPayload []byte) ([]byte, error) {
	in := reflect.New(h.method.Type.In(2).Elem()).Interface()
	err := h.encoding.Unmarshal(inPayload, in)
	if err != nil {
		return nil, fmt.Errorf("encoding unmarshal in data failed:%v", inPayload)
	}
	values := h.method.Func.Call([]reflect.Value{
		reflect.ValueOf(h.receiver),
		reflect.ValueOf(ctx),
		reflect.ValueOf(in),
	})
	out := values[0].Interface()
	outPayload, err := h.encoding.Marshal(out)
	if err != nil {
		return nil, fmt.Errorf("encoding marshal out data failed:%v", out)
	}
	return outPayload, values[1].Interface().(error)
}

func (h *Handler) String() string {
	return fmt.Sprintf("/%s", strings.ToLower(h.method.Name))
}
