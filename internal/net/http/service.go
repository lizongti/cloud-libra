package http

import (
	"context"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/cloudlibraries/libra/internal/boost/coroutine"
	"github.com/cloudlibraries/libra/internal/boost/magic"
	"github.com/cloudlibraries/libra/internal/core/device"
)

type ServiceRequest struct {
	CoroutineID coroutine.ID
	URL         string
	Timeout     time.Duration
	Retry       int
	Proxy       string
	ContentType string
	Form        url.Values
	Body        string
}

type ServiceResponse struct {
	Err        error
	Request    *ServiceRequest
	StatusCode int
	Body       string
}

type Service struct{}

func init() {
	device.Bus().Integrate(&Service{})
}

func (s *Service) HTTP(ctx context.Context, req *ServiceRequest) (resp *ServiceResponse, err error) {
	return s.do(ctx, req, magic.HTTP)
}

func (s *Service) HTTPS(ctx context.Context, req *ServiceRequest) (resp *ServiceResponse, err error) {
	return s.do(ctx, req, magic.HTTPS)
}

func (*Service) do(ctx context.Context, req *ServiceRequest, protocol string) (resp *ServiceResponse, err error) {
	resp = new(ServiceResponse)

	opt := []funcClientOption{WithClientProtocol(protocol), WithClientContext(ctx)}

	if req.Timeout != 0 {
		opt = append(opt, WithClientTimeout(req.Timeout))
	}
	if req.Retry != 0 {
		opt = append(opt, WithClientRetry(req.Retry))
	}
	if req.Proxy != "" {
		opt = append(opt, WithClientProxy(req.Proxy))
	}
	if req.ContentType != "" {
		opt = append(opt, WithClientContentType(req.ContentType))
	}
	if len(req.Form) > 0 {
		opt = append(opt, WithClientForm(req.Form))
	}
	if req.Body != "" {
		opt = append(opt, WithClientRequestBody(func() (io.Reader, error) {
			return strings.NewReader(req.Body), nil
		}))
	}
	c := NewClient(opt...)
	httpResp, body, err := c.Get(req.URL)
	if err != nil {
		return nil, err
	}

	if httpResp != nil {
		resp.StatusCode = httpResp.StatusCode
	}
	resp.Body = string(body)

	return resp, nil
}
