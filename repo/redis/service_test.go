package redis_test

import (
	"context"
	"testing"

	"github.com/aceaura/libra/core/device"
	"github.com/aceaura/libra/core/encoding"
	"github.com/aceaura/libra/core/magic"
	"github.com/aceaura/libra/core/message"
	"github.com/aceaura/libra/core/route"
	"github.com/aceaura/libra/repo/redis"
	"github.com/alicebob/miniredis/v2"
)

func TestService(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	var (
		addr          = s.Addr()
		db            = 0
		ctx           = context.Background()
		encodingStyle = magic.NewChainStyle(magic.SeparatorPeriod, magic.SeparatorUnderscore)
		routeStyle    = magic.NewChainStyle(magic.SeparatorSlash, magic.SeparatorUnderscore)
		e             = encoding.NewChainEncoding(encodingStyle.Chain("json.base64.lazy"), encodingStyle.Chain("lazy.base64.json"))
		r             = route.NewChainRoute(routeStyle.Chain("/client"), routeStyle.Chain("/redis"))
	)
	client := device.NewClient().WithName("Client")
	service := &redis.Service{}
	bus := device.NewRouter().WithBus().WithName("Bus").WithService(service).WithDevice(client)
	t.Logf("\n%s", device.Tree(bus))
	req := &redis.ServiceRequest{
		Addr: addr,
		DB:   db,
		Cmd:  []string{"SET", "test", "100"},
	}
	data, err := e.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	msg := &message.Message{
		Route:    r,
		Encoding: e,
		Data:     data,
	}
	processor := device.NewFuncProcessor(func(ctx context.Context, msg *message.Message) error {
		resp := new(redis.ServiceResponse)
		if err := msg.Encoding.Unmarshal(msg.Data, resp); err != nil {
			return err
		}
		result := resp.Result
		if len(result) == 0 {
			t.Fatal("expected a result with content")
		}
		t.Log(result)
		return nil
	})
	if err = client.Invoke(ctx, msg, processor); err != nil {
		t.Fatalf("unexpected error getting from device: %v", err)
	}
}
