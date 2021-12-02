package redis

import (
	"context"

	"github.com/aceaura/libra/core/device"
)

type RedisActionType int

const (
	Append RedisActionType = iota
	BITCOUNT
	BITOP
	DECR
	DECRBY
	GET
	GETBIT
	GETRANGE
	GETSET
	INCR
	INCRBY
	INCRBYFLOAT
	MGET
	MSET
	MSETNX
	PSETEX
	SET
	SETBIT
	SETEX
	SETNX
	SETRANGE
	STRLEN
)

var redisActionName = map[RedisActionType]string{
	Append:      "APPEND",
	BITCOUNT:    "BITCOUNT",
	BITOP:       "BITOP",
	DECR:        "DECR",
	DECRBY:      "DECRBY",
	GET:         "GET",
	GETBIT:      "GETBIT",
	GETRANGE:    "GETRANGE",
	GETSET:      "GETSET",
	INCR:        "INCR",
	INCRBY:      "INCRBY",
	INCRBYFLOAT: "INCRBYFLOAT",
	MGET:        "MGET",
	MSET:        "MSET",
	MSETNX:      "MSETNX",
	PSETEX:      "PSETEX",
	SET:         "SET",
	SETBIT:      "SETBIT",
	SETEX:       "SETEX",
	SETNX:       "SETNX",
	SETRANGE:    "SETRANGE",
	STRLEN:      "STRLEN",
}

type ServiceRequest struct {
	Addr string
	Cmd  string
}

type ServiceResponse struct {
}

type Service struct{}

func init() {
	device.Bus().WithService(&Service{})
}

func (s *Service) Redis(ctx context.Context, req *ServiceRequest) (resp *ServiceResponse, err error) {
	return nil, nil
}

func (s *Service) RedisPipeline(ctx context.Context, req *ServiceRequest) (resp *ServiceResponse, err error) {
	return nil, nil
}
