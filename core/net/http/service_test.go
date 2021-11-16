package http_test

import (
	"context"
	"testing"
	"time"

	"github.com/aceaura/libra/core/device"
	"github.com/aceaura/libra/core/encoding"
	"github.com/aceaura/libra/core/magic"
	"github.com/aceaura/libra/core/message"
	"github.com/aceaura/libra/core/net/http"
	"github.com/aceaura/libra/core/route"
)

func TestService(t *testing.T) {
	const (
		url = "https://top.baidu.com/board?platform=pc&sa=pcindex_entry"
	)
	var (
		ctx           = context.Background()
		encodingStyle = magic.NewChainStyle(magic.SeparatorPeriod, magic.SeparatorUnderscore)
		routeStyle    = magic.NewChainStyle(magic.SeparatorSlash, magic.SeparatorUnderscore)
		e             = encoding.NewChainEncoding(encodingStyle.Chain("json.base64.lazy"), encodingStyle.Chain("lazy.base64.json"))
		r             = route.NewChainRoute(routeStyle.Chain("/client"), routeStyle.Chain("/https"))
	)
	client := device.NewClient().WithName("Client")
	service := &http.Service{}
	bus := device.NewRouter().WithBus().WithName("Bus").WithService(service).WithDevice(client)

	t.Logf("\n%s", device.Tree(bus))

	req := &http.ServiceRequest{
		URL:     url,
		Timeout: 10 * time.Second,
		Retry:   3,
	}

	data, err := e.Marshal(req)
	if err != nil {
		t.Error(err)
	}

	msg := &message.Message{
		Route:    r,
		Encoding: e,
		Data:     data,
	}
	processor := device.NewFuncProcessor(func(ctx context.Context, msg *message.Message) error {
		resp := new(http.ServiceResponse)
		if err := msg.Encoding.Unmarshal(msg.Data, resp); err != nil {
			return err
		}
		body := resp.Body
		if len(body) == 0 {
			t.Fatal("expected a body with content")
		}
		t.Log(string(body))
		return nil
	})
	if err = client.Request(ctx, msg, processor); err != nil {
		t.Fatalf("unexpected error getting from device: %v", err)
	}
}
