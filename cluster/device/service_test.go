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

var e = encoding.NewJSON()

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
	if err := e.Unmarshal(data, resp); err != nil {
		return err
	}
	s.logChan <- fmt.Sprintf("%v", resp)
	return nil
}

func TestDevice(t *testing.T) {
	const (
		timeout = 10
		version = "1.0.0"
		codec   = "json"
	)
	logChan := make(chan string)
	client := &Client{
		logChan: logChan,
	}
	component := &Try{
		logChan: logChan,
	}
	service := device.NewService(
		device.ServiceOption.WithEncoding(e),
		device.ServiceOption.WithComponent(component),
	)
	router := device.NewRouter(
		device.RouterOption.WithName(version),
		device.RouterOption.WithDevice(service),
	)
	device.NewBus(
		device.BusOption.WithDevice(router),
		device.BusOption.WithDevice(client),
	)
	ctx := context.Background()
	route := device.NewRoute().WithSrc(
		"/client", magic.SeparatorSlash, magic.SeparatorUnderscore,
	).WithDst(
		"/1.0.0/try/echo", magic.SeparatorSlash, magic.SeparatorUnderscore,
	)

	reqData, err := encoding.Marshal(e, &Ping{
		Text: "libra: Hello, world!",
	})
	if err != nil {
		t.Fatalf("unexpected error getting from encoding: %v", err)
	}

	if err = client.Process(ctx, *route, reqData); err != nil {
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

func TestContext(t *testing.T) {

}
