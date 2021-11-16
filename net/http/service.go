package http

import (
	"context"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/aceaura/libra/magic"
)

type ServiceRequest struct {
	URL         string
	Timeout     time.Duration
	Retry       int
	Proxy       string
	ContentType string
	Form        url.Values
	Body        string
}

type ServiceResponse struct {
	StatusCode int
	Body       string
}

type Service struct{}

func (d *Service) HTTP(ctx context.Context, req *ServiceRequest) (resp *ServiceResponse, err error) {
	return d.do(ctx, req, magic.HTTP)
}

func (d *Service) HTTPS(ctx context.Context, req *ServiceRequest) (resp *ServiceResponse, err error) {
	return d.do(ctx, req, magic.HTTPS)
}

func (*Service) do(ctx context.Context, req *ServiceRequest, protocol string) (*ServiceResponse, error) {
	c := NewClient().WithProtocol(protocol).WithContext(ctx)
	if req.Timeout != 0 {
		c.WithTimeout(req.Timeout)
	}
	if req.Retry != 0 {
		c.WithRetry(req.Retry)
	}
	if req.Proxy != "" {
		c.WithProxy(req.Proxy)
	}
	if req.ContentType != "" {
		c.WithContentType(req.ContentType)
	}
	if len(req.Form) > 0 {
		c.WithForm(req.Form)
	}
	if req.Body != "" {
		c.WithRequestBody(func() (io.Reader, error) {
			return strings.NewReader(req.Body), nil
		})
	}
	httpResp, body, err := c.Get(req.URL)
	if err != nil {
		return nil, err
	}
	var resp = new(ServiceResponse)
	if httpResp != nil {
		resp.StatusCode = httpResp.StatusCode
	}
	resp.Body = string(body)

	return resp, nil
}
