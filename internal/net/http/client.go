package http

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const (
	GET  = "GET"
	POST = "POST"
	HEAD = "HEAD"
)

type (
	Request        = http.Request
	Response       = http.Response
	ResponseWriter = http.ResponseWriter
)

type Client struct {
	opts clientOptions
}

func NewClient(opt ...funcClientOption) *Client {
	opts := defaultClientOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	c := &Client{
		opts: opts,
	}

	return c
}

func Get(url string, opts ...funcClientOption) (*http.Response, []byte, error) {
	return NewClient(opts...).Get(url)
}

func (c *Client) Get(url string) (*http.Response, []byte, error) {
	return c.Do(GET, url)
}

func Post(url string, opts ...funcClientOption) (*http.Response, []byte, error) {
	return NewClient(opts...).Post(url)
}

func (c *Client) Post(url string) (*http.Response, []byte, error) {
	return c.Do(POST, url)
}

func Head(url string, opts ...funcClientOption) (*http.Response, []byte, error) {
	return NewClient(opts...).Head(url)
}

func (c *Client) Head(url string) (*http.Response, []byte, error) {
	return c.Do(HEAD, url)
}

func Do(method string, url string, opts ...funcClientOption) (*http.Response, []byte, error) {
	return NewClient(opts...).Do(method, url)
}

func (c *Client) Do(method string, url string) (resp *http.Response, body []byte, err error) {
	if c.opts.safety {
		defer func() {
			if v := recover(); v != nil {
				err = fmt.Errorf("%v", v)
			}
		}()
	}

	return c.do(method, url)
}

func (c *Client) do(method string, httpURL string) (*http.Response, []byte, error) {
	submatch := regexp.MustCompile("(https?://)?(.+)").FindStringSubmatch(httpURL)
	if submatch[1] == "" {
		httpURL = fmt.Sprintf("%s://%s", c.opts.protocol, submatch[0])
	}

	if len(c.opts.form) > 0 {
		httpURL = fmt.Sprintf("%s?%s", httpURL, c.opts.form.Encode())
	}

	client := &http.Client{
		Timeout: c.opts.timeout,
	}

	if c.opts.proxy != "" {
		proxyURL, err := url.Parse(c.opts.proxy)
		if err != nil {
			return nil, nil, err
		}

		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	return c.requestWithRetry(func() (*http.Response, error) {
		var err error
		var body io.Reader
		if c.opts.requestBodyFunc != nil {
			body, err = c.opts.requestBodyFunc()
			if err != nil {
				return nil, err
			}
		}
		req, err := http.NewRequest(method, httpURL, body)
		if err != nil {
			return nil, err
		}
		if c.opts.ctx != nil {
			req = req.WithContext(c.opts.ctx)
		}
		req.Header.Set("Content-Type", c.opts.contentType)

		return client.Do(req)
	})
}

func (c *Client) requestWithRetry(f func() (*http.Response, error)) (resp *http.Response, body []byte, err error) {
	for times := 1; times <= c.opts.retry; times++ {
		resp, body, err = c.request(f)
		if err == nil && resp.StatusCode == 200 {
			break
		}
	}
	return resp, body, err
}

func (c *Client) request(f func() (*http.Response, error)) (resp *http.Response, body []byte, err error) {
	resp, err = f()
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return nil, nil, err
	}
	if c.opts.responseBodyFunc != nil {
		err = c.opts.responseBodyFunc(resp.Body)
		if err != nil {
			return nil, nil, err
		}
	} else {
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, err
		}
	}

	return resp, body, nil
}

func (c *Client) prepareClient() error {

	return nil
}

type clientOptions struct {
	protocol         string
	contentType      string
	retry            int
	form             url.Values
	requestBodyFunc  func() (io.Reader, error)
	responseBodyFunc func(io.Reader) error
	safety           bool
	timeout          time.Duration
	proxy            string
	ctx              context.Context
}

var defaultClientOptions = clientOptions{
	protocol:         "http",
	contentType:      "text/plain",
	retry:            1,
	form:             nil,
	requestBodyFunc:  nil,
	responseBodyFunc: nil,
	safety:           false,
	timeout:          0,
	proxy:            "",
	ctx:              nil,
}

type ApplyClientOption interface {
	apply(*clientOptions)
}

type funcClientOption func(*clientOptions)

func (f funcClientOption) apply(opt *clientOptions) {
	f(opt)
}

func WithClientProtocol(protocol string) funcClientOption {
	return func(c *clientOptions) {
		c.protocol = protocol
	}
}

func WithClientTimeout(timeout time.Duration) funcClientOption {
	return func(c *clientOptions) {
		c.timeout = timeout
	}
}

func WithClientRetry(retry int) funcClientOption {
	return func(co *clientOptions) {
		co.retry = retry
	}
}

func WithClientProxy(proxy string) funcClientOption {
	return func(c *clientOptions) {
		c.proxy = proxy
	}
}

func WithClientContentType(contentType string) funcClientOption {
	return func(c *clientOptions) {
		c.contentType = contentType
	}
}

func WithClientForm(form url.Values) funcClientOption {
	return func(c *clientOptions) {
		c.form = form
	}
}

func WithClientRequestBody(requestBodyFunc func() (io.Reader, error)) funcClientOption {
	return func(c *clientOptions) {
		c.requestBodyFunc = requestBodyFunc
	}
}

func WithClientResponseBodyFunc(responseBodyFunc func(io.Reader) error) funcClientOption {
	return func(co *clientOptions) {
		co.responseBodyFunc = responseBodyFunc
	}
}

func WithClientSafety() funcClientOption {
	return func(c *clientOptions) {
		c.safety = true
	}
}

func WithClientContext(ctx context.Context) funcClientOption {
	return func(c *clientOptions) {
		c.ctx = ctx
	}
}
