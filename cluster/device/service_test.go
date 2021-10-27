package device_test

import (
	"context"
	"testing"

	"github.com/aceaura/libra/cluster/device"
	"github.com/aceaura/libra/encoding"
	"github.com/aceaura/libra/magic"
)

type TestService struct {
	device.Service
}

type TestRequest struct {
	Index int
}
type TestResponse struct {
	Index int
}

func (*TestService) TestHandler(_ context.Context, req *TestRequest) (resp *TestResponse, err error) {
	resp = &TestResponse{Index: req.Index}
	return
}

func TestDevice(t *testing.T) {
	service := new(TestService)
	bus := device.NewBus(
		device.BusOption.WithDevice(service),
	)
	bus.Serve()
	ctx := context.Background()
	route := device.NewRoute().WithSrc(
		"bus/run_test/tmp", magic.SeparatorSlash, magic.SeparatorUnderscore,
	).WithDst(
		"bus/test_service/test_handler", magic.SeparatorSlash, magic.SeparatorUnderscore,
	).Build()

	reqData, err := encoding.Marshal(encoding.JSON(), &TestRequest{
		Index: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error getting from encoding: %v", err)
	}

	if err = bus.Process(ctx, route, reqData); err != nil {
		t.Fatalf("unexpected error getting from device: %v", err)
	}
}
