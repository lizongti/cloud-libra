package redis_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/cloudlibraries/libra/internal/boost/magic"
	"github.com/cloudlibraries/libra/internal/core/device"
	"github.com/cloudlibraries/libra/internal/core/encoding"
	"github.com/cloudlibraries/libra/internal/core/message"
	"github.com/cloudlibraries/libra/internal/core/route"
	"github.com/cloudlibraries/libra/internal/repo/redis"
)

func TestService(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	var (
		url = fmt.Sprintf("redis://:@%s/%d", s.Addr(), 0)
		ctx = context.Background()
		e   = encoding.NewChainEncoding(magic.UnixChain("json.base64.lazy"), magic.UnixChain("lazy.base64.json"))
		r   = route.NewChainRoute(magic.GoogleChain("/client"), magic.GoogleChain("/redis/command"))
	)
	client := device.NewClient("Client")
	service := &redis.Service{}
	redisRouter := device.NewRouter("Redis").Integrate(service)
	bus := device.NewBus().Integrate(redisRouter, client)
	t.Logf("\n%s", device.Tree(bus))
	req := &redis.CommandRequest{
		URL: url,
		Cmd: []string{"SET", "test", "100"},
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
		resp := new(redis.CommandResponse)
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
