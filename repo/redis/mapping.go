package redis

import (
	"context"
	"errors"
	"fmt"

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
	ErrUnsupportedType       = errors.New("unsupported type for mapping")
)

type String struct {
	Key   string
	Value string
}

type Hash struct {
	Key   string
	Value map[interface{}]interface{}
}

type List struct {
	Key   string
	Value []interface{}
}

type Set struct {
	Key   string
	Value []interface{}
}

type SortedSet struct {
	Key   string
	Value map[interface{}]interface{}
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

func (m *Mapping) Replace(v interface{}) error {
	switch v := v.(type) {
	case String:
		return m.replaceString(&v)
	case *String:
		return m.replaceString(v)
	case Hash:
		return m.replaceHash(&v)
	case *Hash:
		return m.replaceHash(v)
	case List:
		return m.replaceList(&v)
	case *List:
		return m.replaceList(v)
	case Set:
		return m.replaceSet(&v)
	case *Set:
		return m.replaceSet(v)
	case SortedSet:
		return m.replaceSortedSet(&v)
	case *SortedSet:
		return m.replaceSortedSet(v)
	default:
		return ErrUnsupportedType
	}
}

func (m *Mapping) replaceString(s *String) error {
	cmd := []string{"SET", s.Key, s.Value}

	result, err := m.invoke(cmd)
	if err != nil {
		return err
	}

	if len(result) != 1 {
		return ErrResultLengthNotValid
	}
	if result[0] != magic.OK {
		return ErrResultContentNotValid
	}
	return nil
}

func (m *Mapping) replaceHash(hash *Hash) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
		}
	}()

	cmd := []string{"DEL", hash.Key}
	result, err := m.invoke(cmd)
	if err != nil {
		return err
	}
	if len(result) != 1 {
		return ErrResultLengthNotValid
	}
	if _, err := cast.ToIntE(result[0]); err != nil {
		return ErrResultContentNotValid
	}

	cmd = make([]string, 0, len(hash.Value)*2+2)
	cmd = append(cmd, "HMSET", hash.Key)
	for k, v := range hash.Value {
		cmd = append(cmd, cast.ToString(k), cast.ToString(v))
	}

	result, err = m.invoke(cmd)
	if err != nil {
		return err
	}

	if len(result) != 1 {
		return ErrResultLengthNotValid
	}
	if result[0] != magic.OK {
		return ErrResultContentNotValid
	}

	return nil
}

func (m *Mapping) replaceList(list *List) error {
	cmd := []string{"DEL", list.Key}
	result, err := m.invoke(cmd)
	if err != nil {
		return err
	}
	if len(result) != 1 {
		return ErrResultLengthNotValid
	}
	if _, err := cast.ToIntE(result[0]); err != nil {
		return ErrResultContentNotValid
	}

	cmd = make([]string, 0, len(list.Value)+2)
	cmd = append(cmd, "LPUSH", list.Key)
	for _, v := range list.Value {
		cmd = append(cmd, cast.ToString(v))
	}
	result, err = m.invoke(cmd)
	if err != nil {
		return err
	}
	if len(result) != 1 {
		return ErrResultLengthNotValid
	}
	if result[0] != magic.OK {
		return ErrResultContentNotValid
	}

	return nil
}

func (m *Mapping) replaceSet(set *Set) error {
	cmd := []string{"DEL", set.Key}
	result, err := m.invoke(cmd)
	if err != nil {
		return err
	}
	if len(result) != 1 {
		return ErrResultLengthNotValid
	}
	if _, err := cast.ToIntE(result[0]); err != nil {
		return ErrResultContentNotValid
	}

	cmd = make([]string, 0, len(set.Value)+2)
	cmd = append(cmd, "SADD", set.Key)
	for _, v := range set.Value {
		cmd = append(cmd, cast.ToString(v))
	}
	result, err = m.invoke(cmd)
	if err != nil {
		return err
	}
	if len(result) != 1 {
		return ErrResultLengthNotValid
	}
	if _, err := cast.ToIntE(result[0]); err != nil {
		return ErrResultContentNotValid
	}
	return nil
}

func (m *Mapping) replaceSortedSet(sortedSet *SortedSet) error {
	cmd := []string{"DEL", sortedSet.Key}
	result, err := m.invoke(cmd)
	if err != nil {
		return err
	}
	if len(result) != 1 {
		return ErrResultLengthNotValid
	}
	if _, err := cast.ToIntE(result[0]); err != nil {
		return ErrResultContentNotValid
	}

	cmd = make([]string, 0, len(sortedSet.Value)*2+2)
	cmd = append(cmd, "ZADD", sortedSet.Key)
	for k, v := range sortedSet.Value {
		cmd = append(cmd, cast.ToString(v), cast.ToString(k))
	}
	result, err = m.invoke(cmd)
	if err != nil {
		return err
	}
	if len(result) != 1 {
		return ErrResultLengthNotValid
	}
	if _, err := cast.ToIntE(result[0]); err != nil {
		return ErrResultContentNotValid
	}
	return nil
}

func (m *Mapping) invoke(cmd []string) (result []string, err error) {
	if err := m.Client.Invoke(m.opts.context, &message.Message{
		Route:    route.NewChainRoute(device.Addr(m), magic.GoogleChain("/redis")),
		Encoding: encoding.NewJSON(),
		Data: encoding.Encode(encoding.NewJSON(), &ServiceRequest{
			URL: m.opts.url,
			Cmd: cmd,
		}),
	}, device.NewFuncProcessor(func(ctx context.Context, msg *message.Message) error {
		resp := new(ServiceResponse)
		encoding.Decode(msg.Encoding, msg.Data, resp)
		result = resp.Result
		return nil
	})); err != nil {
		return nil, err
	}
	return result, nil
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
