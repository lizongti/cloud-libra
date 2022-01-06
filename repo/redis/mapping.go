package redis

import (
	"context"
	"errors"

	"github.com/aceaura/libra/core/cast"
	"github.com/aceaura/libra/core/device"
	"github.com/aceaura/libra/core/encoding"
	"github.com/aceaura/libra/core/magic"
	"github.com/aceaura/libra/core/message"
	"github.com/aceaura/libra/core/route"
)

var (
	ErrResultLengthNotValid  = errors.New("result length is not valid")
	ErrResultContentNotValid = errors.New("result content is not valid")
)

type Hash struct {
	key   string
	value map[string]interface{}
}

type Mapping struct {
	*device.Client
	opts mappingOptions
}

func NewMapping(opt ...ApplyMappingOption) *Mapping {
	opts := defaultMappingOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	return &Mapping{
		Client: device.NewClient(),
		opts:   opts,
	}
}

func (m *Mapping) String() string {
	return m.opts.name
}

func (m *Mapping) storeHash(hash *Hash) error {
	var (
		e = encoding.NewChainEncoding(magic.UnixChain("json.base64.lazy"), magic.UnixChain("lazy.base64.json"))
		r = route.NewChainRoute(magic.GoogleChain("/client"), magic.GoogleChain("/redis"))
	)
	cmd := make([]string, 0, len(hash.value)*2+2)
	cmd = append(cmd, "HMSET", hash.key)
	for k, v := range hash.value {
		cmd = append(cmd, k, cast.ToString(v))
	}
	req := &ServiceRequest{
		URL: m.opts.url,
		Cmd: cmd,
	}
	data, err := e.Marshal(req)
	if err != nil {
		return err
	}
	msg := &message.Message{
		Route:    r,
		Encoding: e,
		Data:     data,
	}
	processor := device.NewFuncProcessor(func(ctx context.Context, msg *message.Message) error {
		resp := new(ServiceResponse)
		if err := msg.Encoding.Unmarshal(msg.Data, resp); err != nil {
			return err
		}
		result := resp.Result
		if len(result) != 1 {
			return ErrResultLengthNotValid
		}
		if result[0] != magic.OK {
			return ErrResultContentNotValid
		}
		return nil
	})
	if err = m.Client.Invoke(m.opts.context, msg, processor); err != nil {
		return err
	}
}

type mappingOptions struct {
	url     string
	name    string
	context context.Context
}

var defaultMappingOptions = mappingOptions{
	url:     "redis://localhost:6379/0",
	name:    "",
	context: context.Background(),
}

type ApplyMappingOption interface {
	apply(*mappingOptions)
}

type funcMappingOption func(*mappingOptions)

func (fmo funcMappingOption) apply(mo *mappingOptions) {
	fmo(mo)
}

type mappingOption int

var MappingOption mappingOption

func (mappingOption) URL(url string) funcMappingOption {
	return func(co *mappingOptions) {
		co.url = url
	}
}

func (c *Mapping) WithURL(url string) *Mapping {
	MappingOption.URL(url).apply(&c.opts)
	return c
}

func (mappingOption) Name(name string) funcMappingOption {
	return func(co *mappingOptions) {
		co.name = name
	}
}

func (c *Mapping) WithName(name string) *Mapping {
	MappingOption.Name(name).apply(&c.opts)
	return c
}

func (mappingOption) Context(context context.Context) funcMappingOption {
	return func(co *mappingOptions) {
		co.context = context
	}
}

func (c *Mapping) WithContext(context context.Context) *Mapping {
	MappingOption.Context(context).apply(&c.opts)
	return c
}
