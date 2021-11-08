package device_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aceaura/libra/core/component"
	"github.com/aceaura/libra/core/device"
	"github.com/aceaura/libra/core/encoding"
	"github.com/aceaura/libra/core/message"
	"github.com/aceaura/libra/core/route"
	"github.com/aceaura/libra/magic"
)

var e1 = encoding.NewJSON()
var e2 = encoding.NewBase64()

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

func (t *Try) EchoBytes(_ context.Context, req []byte) (resp []byte, err error) {
	t.logChan <- fmt.Sprintf("%v", string(req))
	resp = req
	fmt.Println(string(req))
	return resp, err
}

type Client1 struct {
	*device.Base
	logChan chan<- string
}

func (s *Client1) Process(ctx context.Context, msg *message.Message) error {
	if msg.Route.Assembling() {
		return s.Gateway().Process(ctx, msg)
	}
	resp := &Pong{}
	if err := e1.Unmarshal(msg.Data, resp); err != nil {
		return err
	}
	s.logChan <- fmt.Sprintf("%v", resp)
	return nil
}

type Client2 struct {
	*device.Base
	logChan chan<- string
}

func (s *Client2) Process(ctx context.Context, msg *message.Message) error {
	if msg.Route.Assembling() {
		return s.Gateway().Process(ctx, msg)
	}
	bytes := &encoding.Bytes{}
	if err := e2.Unmarshal(msg.Data, bytes); err != nil {
		return err
	}
	s.logChan <- fmt.Sprintf("%v", string(bytes.Data))
	return nil
}

func TestService1(t *testing.T) {
	const (
		timeout     = 10
		version     = "1.0.0"
		logChanSize = 2
		msgID       = 0
	)
	logChan := make(chan string, logChanSize)
	client := &Client1{
		Base:    device.NewBase(),
		logChan: logChan,
	}
	component := &Try{
		logChan: logChan,
	}
	service := device.NewService(
		device.ServiceOption.WithEncoding(e2),
		device.ServiceOption.WithComponent(component),
	)
	router := device.NewRouter(
		device.RouterOption.WithName(version),
		device.RouterOption.WithDevice(service),
	)
	bus := device.NewBus(
		device.BusOption.WithDevice(client),
		device.BusOption.WithDevice(router),
	)

	t.Logf("\n%s", device.Tree(bus))

	ctx := context.Background()
	r := route.NewRoute().WithSrc(
		"/anonymous", magic.SeparatorSlash, magic.SeparatorUnderscore,
	).WithDst(
		"/1.0.0/try/echo", magic.SeparatorSlash, magic.SeparatorUnderscore,
	)

	reqData, err := encoding.Marshal(e1, &Ping{
		Text: "libra: Hello, world!",
	})
	if err != nil {
		t.Fatalf("unexpected error getting from encoding: %v", err)
	}
	msg := &message.Message{
		ID:       msgID,
		Route:    *r,
		Encoding: e1,
		Data:     reqData,
	}

	if err = client.Process(ctx, msg); err != nil {
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

func TestService2(t *testing.T) {
	const (
		timeout     = 10
		version     = "1.0.0"
		logChanSize = 2
		msgID       = 0
	)
	logChan := make(chan string, logChanSize)
	client := &Client2{
		Base:    device.NewBase(),
		logChan: logChan,
	}
	component := &Try{
		logChan: logChan,
	}
	service := device.NewService(
		device.ServiceOption.WithEncoding(e2),
		device.ServiceOption.WithComponent(component),
	)
	router := device.NewRouter(
		device.RouterOption.WithName(version),
		device.RouterOption.WithDevice(service),
	)
	bus := device.NewBus(
		device.BusOption.WithDevice(client),
		device.BusOption.WithDevice(router),
	)

	t.Logf("\n%s", device.Tree(bus))

	ctx := context.Background()
	r := route.NewRoute().WithSrc(
		"/anonymous", magic.SeparatorSlash, magic.SeparatorUnderscore,
	).WithDst(
		"/1.0.0/try/echo_bytes", magic.SeparatorSlash, magic.SeparatorUnderscore,
	)

	reqData, err := encoding.Marshal(e2, []byte("libra: Hello, world!"))
	if err != nil {
		t.Fatalf("unexpected error getting from encoding: %v", err)
	}
	msg := &message.Message{
		ID:       msgID,
		Route:    *r,
		Encoding: encoding.Empty(),
		Data:     reqData,
	}

	if err = client.Process(ctx, msg); err != nil {
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
