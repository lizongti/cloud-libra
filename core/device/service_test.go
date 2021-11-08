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
	*device.Base
	logChan chan<- string
}

func (s *Client) Process(ctx context.Context, msg *message.Message) error {
	if msg.State() == message.MessageStateAssembling {
		return s.Gateway().Process(ctx, msg)
	}
	resp := &Pong{}
	if err := e.Unmarshal(msg.Data(), resp); err != nil {
		return err
	}
	s.logChan <- fmt.Sprintf("%v", resp)
	return nil
}

func TestDevice(t *testing.T) {
	const (
		timeout     = 10
		version     = "1.0.0"
		logChanSize = 2
	)
	logChan := make(chan string, logChanSize)
	client := &Client{
		Base:    device.NewBase(),
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
	bus := device.NewBus(
		device.BusOption.WithDevice(client),
		device.BusOption.WithDevice(router),
	)

	t.Logf("\n%s", device.Tree(bus))

	ctx := context.Background()
	route := route.NewRoute().WithSrc(
		"/anonymous", magic.SeparatorSlash, magic.SeparatorUnderscore,
	).WithDst(
		"/1.0.0/try/echo", magic.SeparatorSlash, magic.SeparatorUnderscore,
	)

	reqData, err := encoding.Marshal(e, &Ping{
		Text: "libra: Hello, world!",
	})
	if err != nil {
		t.Fatalf("unexpected error getting from encoding: %v", err)
	}
	msg := message.NewMessage(0, *route, encoding.Empty(), reqData)

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
