package redis

import (
	"context"

	"github.com/aceaura/libra/core/device"
)

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
