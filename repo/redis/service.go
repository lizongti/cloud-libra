package redis

import (
	"context"

	"github.com/aceaura/libra/core/cast"
	"github.com/aceaura/libra/core/device"
)

type ServiceRequest struct {
	Addr string
	DB   int
	Cmd  []string
}

type ServiceResponse struct {
	Err     error
	Request *ServiceRequest
	Result  []string
}

type Service struct{}

func init() {
	device.Bus().WithService(&Service{})
}

func (s *Service) Redis(ctx context.Context, req *ServiceRequest) (resp *ServiceResponse, err error) {
	resp = new(ServiceResponse)
	c := NewClient().WithAddr(req.Addr).WithDB(req.DB).WithContext(ctx)
	result, err := c.Command(cast.ToSlice(req.Cmd)...)
	if err != nil {
		return nil, err
	}
	resp.Result = result
	return resp, nil
}

func (s *Service) RedisPipeline(ctx context.Context, req *ServiceRequest) (resp *ServiceResponse, err error) {
	return nil, nil
}
