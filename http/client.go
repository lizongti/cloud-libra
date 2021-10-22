package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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
	opts               []clientOpt
	protocol           string
	contentType        string
	timeout            time.Duration
	retry              int
	proxy              string
	form               url.Values
	body               io.Reader
	responseBodyReader func(io.Reader) error
	safety             bool
	client             http.Client
}

func NewClient(opts ...clientOpt) *Client {
	return &Client{opts: opts}
}

func Get(url string, opts ...clientOpt) (*http.Response, []byte, error) {
	return NewClient(opts...).Get(url)
}

func (c *Client) Get(url string) (*http.Response, []byte, error) {
	return c.Do(GET, url)
}

func Post(url string, opts ...clientOpt) (*http.Response, []byte, error) {
	return NewClient(opts...).Post(url)
}

func (c *Client) Post(url string) (*http.Response, []byte, error) {
	return c.Do(POST, url)
}

func Head(url string, opts ...clientOpt) (*http.Response, []byte, error) {
	return NewClient(opts...).Head(url)
}

func (c *Client) Head(url string) (*http.Response, []byte, error) {
	return c.Do(HEAD, url)
}

func Do(method string, url string, opts ...clientOpt) (*http.Response, []byte, error) {
	return NewClient(opts...).Do(method, url)
}

func (c *Client) Do(method string, url string) (resp *http.Response, body []byte, err error) {
	c.init()
	if c.safety {
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("%v", e)
			}
		}()
	}
	return c.do(method, url)
}

func (c *Client) init() error {
	c.form = make(url.Values)

	for _, opt := range c.opts {
		opt(c)
	}

	if c.proxy != "" {
		fixedURL, err := url.Parse(c.proxy)
		if err != nil {
			return err
		}
		c.client.Transport = &http.Transport{Proxy: http.ProxyURL(fixedURL)}
	}
	if c.timeout > 0 {
		c.client.Timeout = c.timeout
	}
	if c.retry < 1 {
		c.retry = 1
	}
	if c.contentType == "" {
		c.contentType = "text/plain"
	}
	return nil
}

func (c *Client) do(method string, url string) (*http.Response, []byte, error) {
	return c.requestWithRetry(func() (*http.Response, error) {
		submatch := regexp.MustCompile("(https?://)?(.+)").FindStringSubmatch(url)
		if submatch[1] == "" {
			if c.protocol == "" {
				c.protocol = "http"
			}
			url = fmt.Sprintf("%s://%s", c.protocol, submatch[0])
		}

		if len(c.form) > 0 {
			url = fmt.Sprintf("%s?%s", url, c.form.Encode())
		}
		req, err := http.NewRequest(method, url, c.body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", c.contentType)
		return c.client.Do(req)
	})
}

func (c *Client) requestWithRetry(f func() (*http.Response, error)) (resp *http.Response, body []byte, err error) {
	for times := 1; times <= c.retry; times++ {
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
	if c.responseBodyReader != nil {
		err = c.responseBodyReader(resp.Body)
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

type clientOpt func(*Client)
type clientOption struct{}

var ClientOption clientOption

func (clientOption) WithProtocol(protocol string) clientOpt {
	return func(c *Client) {
		c.protocol = protocol
	}
}

func (c *Client) WithProtocol(protocol string) *Client {
	c.opts = append(c.opts, ClientOption.WithProtocol(protocol))
	return c
}

func (clientOption) WithTimeout(timeout time.Duration) clientOpt {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.opts = append(c.opts, ClientOption.WithTimeout(timeout))
	return c
}

func (clientOption) WithRetry(retry int) clientOpt {
	return func(c *Client) {
		c.retry = retry
	}
}

func (c *Client) WithRetry(retry int) *Client {
	c.opts = append(c.opts, ClientOption.WithRetry(retry))
	return c
}

func (clientOption) WithProxy(proxy string) clientOpt {
	return func(c *Client) {
		c.proxy = proxy
	}
}

func (c *Client) WithProxy(proxy string) *Client {
	c.opts = append(c.opts, ClientOption.WithProxy(proxy))
	return c
}

func (clientOption) WithContentType(contentType string) clientOpt {
	return func(c *Client) {
		c.contentType = contentType
	}
}

func (c *Client) WithContentType(contentType string) *Client {
	c.opts = append(c.opts, ClientOption.WithContentType(contentType))
	return c
}

func (clientOption) WithForm(form url.Values) clientOpt {
	return func(c *Client) {
		c.form = form
	}
}

func (c *Client) WithForm(form url.Values) *Client {
	c.opts = append(c.opts, ClientOption.WithForm(form))
	return c
}

func (clientOption) WithParam(key string, value string) clientOpt {
	return func(c *Client) {
		c.form.Set(key, value)
	}
}

func (c *Client) WithParam(key string, value string) *Client {
	c.opts = append(c.opts, ClientOption.WithParam(key, value))
	return c
}

func (clientOption) WithBody(body interface{}) clientOpt {
	return func(c *Client) {
		switch body := body.(type) {
		case string:
			c.body = strings.NewReader(body)
		case []byte:
			c.body = bytes.NewReader(body)
		case io.Reader:
			c.body = body
		default:
			c.body = strings.NewReader(fmt.Sprintf("%v", body))
		}
	}
}

func (c *Client) WithBody(body interface{}) *Client {
	c.opts = append(c.opts, ClientOption.WithBody(body))
	return c
}

func (clientOption) WithResponseBodyReader(responseBodyReader func(io.Reader) error) clientOpt {
	return func(c *Client) {
		c.responseBodyReader = responseBodyReader
	}
}

func (c *Client) WithResponseBodyReader(responseBodyReader func(io.Reader) error) *Client {
	c.opts = append(c.opts, ClientOption.WithResponseBodyReader(responseBodyReader))
	return c
}

func (clientOption) WithSafety() clientOpt {
	return func(c *Client) {
		c.safety = true
	}
}

func (c *Client) WithSafety() *Client {
	c.opts = append(c.opts, ClientOption.WithSafety())
	return c
}
