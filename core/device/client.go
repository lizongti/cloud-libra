package device

import (
	"context"

	"github.com/aceaura/libra/core/message"
)

type Client struct {
	*Base
	name    string
	process func(context.Context, *message.Message) error
}

func NewClient(name string, process func(context.Context, *message.Message) error) *Client {
	return &Client{
		Base:    NewBase(),
		process: process,
	}
}

func (r *Client) String() string {
	return r.name
}

func (t *Client) Process(ctx context.Context, msg *message.Message) error {
	return t.process(ctx, msg)
}
