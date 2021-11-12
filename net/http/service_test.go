package http_test

import (
	"context"
	"testing"

	"github.com/aceaura/libra/core/device"
	"github.com/aceaura/libra/core/message"
	"github.com/aceaura/libra/core/route"
	"github.com/aceaura/libra/magic"
	"github.com/aceaura/libra/net/http"
)

func TestService(t *testing.T) {
	s := http.Service{}
	client := device.NewClient(func(ctx context.Context, msg *message.Message) error {

	})
	device.Bus().WithService(s).WithDevice(client)
	ctx := context.Background()
	r := route.New(
		magic.ChainSplashUnderScore("/client"),
		magic.ChainSplashUnderScore("/http"),
	)
	msg := &message.Message{
		ID:    1,
		Route: r,
		// Encoding : "",
		// Data    :
	}

	device.Bus().Process(ctx, msg)

}
