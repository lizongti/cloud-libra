package redis

import (
	"context"

	"github.com/aceaura/libra/boost/cast"
	"github.com/aceaura/libra/core/device"
)

type CommandRequest struct {
	URL string
	Cmd []string
}

type CommandResponse struct {
	Result []string
}

type PipelineRequest struct {
}

type PipelineResponse struct {
}

type Service struct{}

func init() {
	router := device.NewRouter("Redis").Integrate(&Service{})
	device.Bus().Integrate(router)
}

func (s *Service) Command(ctx context.Context, req *CommandRequest) (resp *CommandResponse, err error) {
	resp = new(CommandResponse)
	c := NewClient(req.URL,
		WithClientContext(ctx),
	)
	result, err := c.Command(cast.ToSlice(req.Cmd)...)
	if err != nil {
		return nil, err
	}
	resp.Result = result
	return resp, nil
}

// TODO ADD POOL

func (s *Service) Pipeline(ctx context.Context, req *PipelineRequest) (resp *PipelineResponse, err error) {
	return nil, nil
}
