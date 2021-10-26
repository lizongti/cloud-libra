package device_test

import (
	"context"
	"testing"

	"github.com/aceaura/libra/cluster/device"
)

type TestService struct {
	device.Service
}

type TestHandlerRequest struct {
	Index int
}
type TestHandlerResponse struct {
	Index int
}

func (*TestService) TestHandler(_ context.Context, req *TestHandlerRequest) (resp *TestHandlerResponse, err error) {
	resp = &TestHandlerResponse{Index: req.Index}
	return
}

func TestDevice(t *testing.T) {
	// service := new(TestService)
	// bus := device.NewBus(
	// 	device.BusOption.WithDevice(service),
	// )
	// ctx := context.Background()
	// route := device.NewRoute().WithSrc("")
	// bus.Process(ctx)
}
