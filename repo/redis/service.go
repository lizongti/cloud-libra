package redis

import (
	"context"

	"github.com/aceaura/libra/core/device"
)

type ServiceRequest struct {
}

type ServiceResponse struct {
}

type Service struct{}

func init() {
	device.Bus().WithService(&Service{})
}

func (d *Service) Redis(ctx context.Context, req *ServiceRequest) (resp *ServiceResponse, err error) {
	return nil, nil
}

func (d *Service) RedisPipeline(ctx context.Context, req *ServiceRequest) (resp *ServiceResponse, err error) {
	return nil, nil
}
