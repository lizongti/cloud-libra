package device

import (
	"context"
	"reflect"

	"github.com/aceaura/libra/scheduler"
)

var (
	typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()
	typeOfBytes   = reflect.TypeOf(([]byte)(nil))
)

type Handler struct {
	Device
	method  reflect.Method
	gateway Device
}

func (h *Handler) String() string {
	return h.method.Name
}

func (h *Handler) Gateway(device Device) {
	h.gateway = device
}

func (h *Handler) Process(ctx context.Context, route Route, data []byte) error {
	deviceType := route.deviceType()
	if deviceType == DeviceTypeBus {
		return h.gateway.Process(ctx, route, data)
	} else if deviceType == DeviceTypeHandler {
		return h.localProcess(ctx, route, data)
	}

	return ErrRouteDeadEnd
}

func (h *Handler) localProcess(ctx context.Context, route Route, reqData []byte) error {
	s := h.gateway.(*Service)
	stage := func(t *scheduler.Task) error {
		respData, err := h.Do(ctx, reqData)
		if err != nil {
			return err
		}
		return s.Process(ctx, route.reverse(), respData)
	}
	scheduler.NewTask(
		scheduler.TaskOption.WithContext(ctx),
		scheduler.TaskOption.WithStages(stage),
	).Publish(s.schedulerFunc(ctx))
	return nil
}

func (h *Handler) Do(ctx context.Context, reqData []byte) (respData []byte, err error) {
	s := h.gateway.(*Service)
	mt := h.method.Type
	var req interface{}
	if mt.In(2) == typeOfBytes {
		req = reqData
	} else {
		req = reflect.New(mt.In(2).Elem()).Interface()
		if err = s.codec.Unmarshal(reqData, req); err != nil {
			return nil, err
		}
	}
	context := new(context.Context)
	in := []reflect.Value{reflect.ValueOf(s), reflect.ValueOf(context), reflect.ValueOf(req)}

	out := h.method.Func.Call(in)
	if err = out[1].Interface().(error); err != nil {
		return nil, err
	}
	resp := out[0].Interface().(error)
	if respData, err = s.codec.Marshal(resp); err != nil {
		return nil, err
	}
	return
}
