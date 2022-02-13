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
	"github.com/mohae/deepcopy"
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
	Value map[string]string
}

type List struct {
	Key   string
	Value []string
}

type Set struct {
	Key   string
	Value []string
}

type SortedSet struct {
	Key   string
	Value map[string]string
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

func (m *Mapping) Append(v interface{}) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
		}
	}()

	switch v := v.(type) {
	case Hash:
		return m.appendHash(&v)
	case *Hash:
		return m.appendHash(v)
	default:
		return ErrUnsupportedType
	}
}

func (m *Mapping) Replace(v interface{}) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
		}
	}()

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

func (m *Mapping) Delete(v interface{}) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
		}
	}()

	var key string

	switch v := v.(type) {
	case String:
		key = v.Key
	case *String:
		key = v.Key
	case Hash:
		key = v.Key
	case *Hash:
		key = v.Key
	case List:
		key = v.Key
	case *List:
		key = v.Key
	case Set:
		key = v.Key
	case *Set:
		key = v.Key
	case SortedSet:
		key = v.Key
	case *SortedSet:
		key = v.Key
	default:
		return ErrUnsupportedType
	}

	cmd := []string{"DEL", key}
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
	return nil
}

func (m *Mapping) Select(v interface{}) (_ interface{}, err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
		}
	}()

	switch v := v.(type) {
	case String:
		v = deepcopy.Copy(v).(String)
		return m.selectString(&v)
	case *String:
		return m.selectString(v)
	case Hash:
		v = deepcopy.Copy(v).(Hash)
		return m.selectHash(&v)
	case *Hash:
		return m.selectHash(v)
	case List:
		v = deepcopy.Copy(v).(List)
		return m.selectList(&v)
	case *List:
		return m.selectList(v)
	case Set:
		v = deepcopy.Copy(v).(Set)
		return m.selectSet(&v)
	case *Set:
		return m.selectSet(v)
	case SortedSet:
		v = deepcopy.Copy(v).(SortedSet)
		return m.selectSortedSet(&v)
	case *SortedSet:
		return m.selectSortedSet(v)
	default:
		return v, ErrUnsupportedType
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

func (m *Mapping) appendHash(hash *Hash) error {
	if len(hash.Value) == 0 {
		return nil
	}

	cmd := make([]string, 0, len(hash.Value)*2+2)
	cmd = append(cmd, "HMSET", hash.Key)
	for k, v := range hash.Value {
		cmd = append(cmd, k, v)
	}

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

func (m *Mapping) replaceHash(hash *Hash) error {
	m.Delete(hash)

	if len(hash.Value) == 0 {
		return nil
	}

	cmd := make([]string, 0, len(hash.Value)*2+2)
	cmd = append(cmd, "HMSET", hash.Key)
	for k, v := range hash.Value {
		cmd = append(cmd, k, v)
	}

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

func (m *Mapping) replaceList(list *List) error {
	m.Delete(list)

	if len(list.Value) == 0 {
		return nil
	}

	cmd := make([]string, 0, len(list.Value)+2)
	cmd = append(cmd, "RPUSH", list.Key)
	for _, v := range list.Value {
		cmd = append(cmd, v)
	}
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

	return nil
}

func (m *Mapping) replaceSet(set *Set) error {
	m.Delete(set)

	if len(set.Value) == 0 {
		return nil
	}

	cmd := make([]string, 0, len(set.Value)+2)
	cmd = append(cmd, "SADD", set.Key)
	for _, v := range set.Value {
		cmd = append(cmd, v)
	}
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
	return nil
}

func (m *Mapping) replaceSortedSet(sortedSet *SortedSet) error {
	m.Delete(sortedSet)

	if len(sortedSet.Value) == 0 {
		return nil
	}

	cmd := make([]string, 0, len(sortedSet.Value)*2+2)
	cmd = append(cmd, "ZADD", sortedSet.Key)
	for k, v := range sortedSet.Value {
		cmd = append(cmd, v, k)
	}
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
	return nil
}

func (m *Mapping) selectString(s *String) (*String, error) {
	cmd := []string{"GET", s.Key}
	result, err := m.invoke(cmd)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, ErrResultLengthNotValid
	}

	s.Value = result[0]
	return s, nil
}

func (m *Mapping) selectHash(hash *Hash) (*Hash, error) {
	cmd := []string{"HGETALL", hash.Key}
	result, err := m.invoke(cmd)
	if err != nil {
		return nil, err
	}
	if len(result)%2 != 0 {
		return nil, ErrResultLengthNotValid
	}

	resultMap := map[string]string{}
	for i, v := range result {
		if i%2 == 1 {
			resultMap[result[i-1]] = v
		}
	}

	hash.Value = resultMap
	return hash, nil
}

func (m *Mapping) selectList(list *List) (*List, error) {
	cmd := []string{"LRANGE", list.Key, "0", "-1"}
	result, err := m.invoke(cmd)
	if err != nil {
		return nil, err
	}

	list.Value = result
	return list, nil
}

func (m *Mapping) selectSet(set *Set) (*Set, error) {
	cmd := []string{"SMEMBERS", set.Key}
	result, err := m.invoke(cmd)
	if err != nil {
		return nil, err
	}

	set.Value = result
	return set, nil
}

func (m *Mapping) selectSortedSet(sortedSet *SortedSet) (*SortedSet, error) {
	cmd := []string{"ZRANGE", sortedSet.Key, "0", "-1", "WITHSCORES"}
	result, err := m.invoke(cmd)
	if err != nil {
		return nil, err
	}
	if len(result)%2 != 0 {
		return nil, ErrResultLengthNotValid
	}

	resultMap := map[string]string{}
	for i, v := range result {
		if i%2 == 1 {
			resultMap[result[i-1]] = v
		}
	}

	sortedSet.Value = resultMap
	return sortedSet, nil
}

func (m *Mapping) invoke(cmd []string) (result []string, err error) {
	if err := m.Client.Invoke(m.opts.context, &message.Message{
		Route:    route.NewChainRoute(device.Addr(m), magic.GoogleChain("/redis/command")),
		Encoding: encoding.NewJSON(),
		Data: encoding.Encode(encoding.NewJSON(), &CommandRequest{
			URL: m.opts.url,
			Cmd: cmd,
		}),
	}, device.NewFuncProcessor(func(ctx context.Context, msg *message.Message) error {
		resp := new(CommandResponse)
		encoding.Decode(msg.Encoding, msg.Data, resp)
		if resp.Result == nil {
			result = make([]string, 0)
		} else {
			result = resp.Result
		}
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

func (f funcMappingOption) apply(opt *mappingOptions) {
	f(opt)
}

type mappingOption int

var MappingOption mappingOption

func (mappingOption) URL(url string) funcMappingOption {
	return func(c *mappingOptions) {
		c.url = url
	}
}

func (c *Mapping) WithURL(url string) *Mapping {
	MappingOption.URL(url).apply(&c.opts)
	return c
}

func (mappingOption) Name(name string) funcMappingOption {
	return func(c *mappingOptions) {
		c.name = name
	}
}

func (c *Mapping) WithName(name string) *Mapping {
	MappingOption.Name(name).apply(&c.opts)
	return c
}

func (mappingOption) Context(context context.Context) funcMappingOption {
	return func(c *mappingOptions) {
		c.context = context
	}
}

func (c *Mapping) WithContext(context context.Context) *Mapping {
	MappingOption.Context(context).apply(&c.opts)
	return c
}
