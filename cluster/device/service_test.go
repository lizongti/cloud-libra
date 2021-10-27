package device_test

import (
	"context"
	"testing"

	"github.com/aceaura/libra/cluster/component"
	"github.com/aceaura/libra/cluster/device"
	"github.com/aceaura/libra/encoding"
	"github.com/aceaura/libra/magic"
)

type Try struct {
	component.ComponentBase
}

type Ping struct {
	Text string
}
type Pong struct {
	Text string
}

func (*Try) Echo(_ context.Context, req *Ping) (resp *Pong, err error) {
	resp = &Pong{Text: req.Text}
	return
}

func TestDevice(t *testing.T) {
	service := device.NewService(
		device.ServiceOption.WithEncoding(encoding.JSON()),
		device.ServiceOption.WithComponent(&Try{}),
	)
	bus := device.NewBus(
		device.BusOption.WithDevice(service),
	)
	ctx := context.Background()
	route := device.NewRoute().WithSrc(
		"bus/one/shot", magic.SeparatorSlash, magic.SeparatorUnderscore,
	).WithDst(
		"bus/try/echo", magic.SeparatorSlash, magic.SeparatorUnderscore,
	)

	reqData, err := encoding.Marshal(encoding.JSON(), &Ping{
		Text: "libra: Hello, world!",
	})
	if err != nil {
		t.Fatalf("unexpected error getting from encoding: %v", err)
	}

	if err = bus.Process(ctx, *route, reqData); err != nil {
		t.Fatalf("unexpected error getting from device: %v", err)
	}
}
