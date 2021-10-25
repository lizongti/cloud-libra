package device_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/aceaura/libra/cluster/device"
	"github.com/aceaura/libra/codec"
)

type TestDevice struct {
}

func (*TestDevice) Code(reflect.Method) uint64 {
	return 0
}

func (*TestDevice) Name(reflect.Method) string {
	return ""
}

func (*TestDevice) Codec(reflect.Method) codec.Codec {
	return nil
}

type TestHandlerRequest struct {
	Index int
}
type TestHandlerResponse struct {
	Index int
}

func (*TestDevice) TestHandler(_ context.Context, req *TestHandlerRequest) (resp *TestHandlerResponse, err error) {
	resp = &TestHandlerResponse{Index: req.Index}
	return
}

type TestLinkRouter struct {
	device.Router
	Handlers []*device.Handler
}

func (tr *TestLinkRouter) Register(h *device.Handler) {
	tr.Handlers = append(tr.Handlers, h)
}

func TestLink(t *testing.T) {
	// r := new(router.Router)
	// r.Link(new(TestDevice))
	// if len(r.Handlers) == 0 {
	// 	t.Fatalf("expect a handle in router, got nil")
	// }
}
