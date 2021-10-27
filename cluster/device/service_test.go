package device_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aceaura/libra/cluster/component"
	"github.com/aceaura/libra/cluster/device"
	"github.com/aceaura/libra/encoding"
	"github.com/aceaura/libra/magic"
)

type Try struct {
	component.ComponentBase
	logChan chan<- string
}

type Ping struct {
	Text string
}
type Pong struct {
	Text string
}

func (t *Try) Echo(_ context.Context, req *Ping) (resp *Pong, err error) {
	t.logChan <- fmt.Sprintf("%v", req)
	resp = &Pong{Text: req.Text}
	return resp, err
}

type Client struct {
	logChan chan<- string
	gateway device.Device
}

func (*Client) String() string {
	return "Client"
}

func (s *Client) LinkGateway(device device.Device) {
	s.gateway = device
}

func (s *Client) Process(ctx context.Context, route device.Route, data []byte) error {
	if route.Taking() {
		return s.gateway.Process(ctx, route, data)
	}
	resp := &Pong{}
	if err := encoding.JSON().Unmarshal(data, resp); err != nil {
		return err
	}
	s.logChan <- fmt.Sprintf("%v", resp)
	return nil
}

func TestDevice(t *testing.T) {
	const (
		timeout = 100
	)
	logChan := make(chan string)
	component := &Try{
		logChan: logChan,
	}
	service := device.NewService(
		device.ServiceOption.WithEncoding(encoding.JSON()),
		device.ServiceOption.WithComponent(component),
	)
	client := &Client{
		logChan: logChan,
	}
	bus := device.NewBus(
		device.BusOption.WithDevice(service),
		device.BusOption.WithDevice(client),
	)
	ctx := context.Background()
	route := device.NewRoute().WithSrc(
		"bus/client", magic.SeparatorSlash, magic.SeparatorUnderscore,
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
	var timeoutChan = time.After(time.Duration(timeout) * time.Second)
	var in string
	var out string
	for {
		select {
		case <-timeoutChan:
			t.Fatal("timeout when getting report from task")
		case msg := <-logChan:
			if in == "" {
				in = msg
				break
			}
			out = msg
			if in != out {
				t.Fatal("expecting out msg equals to in msg")
			}
			return
		}
	}
}
