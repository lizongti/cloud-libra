package device

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/aceaura/libra/core/message"
)

var ErrMissingProcessor = errors.New("processor cannot be found by message ID")

type Processor interface {
	Process(context.Context, *message.Message) error
}

type funcProcessor func(context.Context, *message.Message) error

func NewFuncProcessor(f funcProcessor) funcProcessor {
	return f
}

func (fp funcProcessor) Process(ctx context.Context, msg *message.Message) error {
	return fp(ctx, msg)
}

type Client struct {
	*Base
	opts       clientOptions
	msgID      uint64
	processors sync.Map
}

func NewClient() *Client {
	return &Client{
		Base: NewBase(),
	}
}

func (c *Client) String() string {
	return c.opts.name
}

func (c *Client) Process(ctx context.Context, msg *message.Message) error {
	if msg.Route.Assembling() {
		return c.gateway.Process(ctx, msg)
	}
	return c.localProcess(ctx, msg)
}

func (c *Client) localProcess(ctx context.Context, m *message.Message) error {
	v, ok := c.processors.Load(m.ID)
	if !ok {
		return ErrMissingProcessor
	}
	c.processors.Delete(m.ID)
	return v.(Processor).Process(ctx, m)
}

func (c *Client) Request(ctx context.Context, m *message.Message, p Processor) error {
	m.ID = atomic.AddUint64(&c.msgID, 1)
	c.processors.Store(m.ID, p)
	return c.Process(ctx, m)
}

type clientOptions struct {
	name string
}

var defaultClientOptions = clientOptions{
	name: "",
}

type ApplyClientOption interface {
	apply(*clientOptions)
}

type funcClientOption func(*clientOptions)

func (fco funcClientOption) apply(co *clientOptions) {
	fco(co)
}

type clientOption int

var ClientOption clientOption

func (clientOption) Name(name string) funcClientOption {
	return func(c *clientOptions) {
		c.name = name
	}
}

func (c *Client) WithName(name string) *Client {
	ClientOption.Name(name).apply(&c.opts)
	return c
}
