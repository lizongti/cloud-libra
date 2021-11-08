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
	protocol           string
	contentType        string
	retry              int
	form               url.Values
	body               io.Reader
	responseBodyReader func(io.Reader) error
	safety             bool
	client             http.Client
}

func NewClient(opts ...funcClientOption) *Client {
	c := &Client{
		form:        make(url.Values),
		contentType: "text/plain",
		retry:       1,
	}
	for _, opt := range opts {
		opt(c)
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
	if c.safety {
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("%v", e)
			}
		}()
	}
	return c.do(method, url)
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

type funcClientOption func(*Client)
type clientOption struct{}

var ClientOption clientOption

func (clientOption) WithProtocol(protocol string) funcClientOption {
	return func(c *Client) {
		c.WithProtocol(protocol)
	}
}

func (c *Client) WithProtocol(protocol string) *Client {
	c.protocol = protocol
	return c
}

func (clientOption) WithTimeout(timeout time.Duration) funcClientOption {
	return func(c *Client) {
		c.WithTimeout(timeout)
	}
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.client.Timeout = timeout
	return c
}

func (clientOption) WithRetry(retry int) funcClientOption {
	return func(c *Client) {
		c.WithRetry(retry)
	}
}

func (c *Client) WithRetry(retry int) *Client {
	c.retry = retry
	return c
}

func (clientOption) WithProxy(proxy string) funcClientOption {
	return func(c *Client) {
		c.WithProxy(proxy)
	}
}

func (c *Client) WithProxy(proxy string) *Client {
	fixedURL, err := url.Parse(proxy)
	if err != nil {
		panic(err)
	}
	c.client.Transport = &http.Transport{Proxy: http.ProxyURL(fixedURL)}
	return c
}

func (clientOption) WithContentType(contentType string) funcClientOption {
	return func(c *Client) {
		c.WithContentType(contentType)
	}
}

func (c *Client) WithContentType(contentType string) *Client {
	c.contentType = contentType
	return c
}

func (clientOption) WithForm(form url.Values) funcClientOption {
	return func(c *Client) {
		c.WithForm(form)
	}
}

func (c *Client) WithForm(form url.Values) *Client {
	c.form = form
	return c
}

func (clientOption) WithParam(key string, value string) funcClientOption {
	return func(c *Client) {
		c.WithParam(key, value)
	}
}

func (c *Client) WithParam(key string, value string) *Client {
	c.form.Set(key, value)
	return c
}

func (clientOption) WithBody(body interface{}) funcClientOption {
	return func(c *Client) {
		c.WithBody(body)
	}
}

func (c *Client) WithBody(body interface{}) *Client {
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
	return c
}

func (clientOption) WithResponseBodyReader(responseBodyReader func(io.Reader) error) funcClientOption {
	return func(c *Client) {
		c.WithResponseBodyReader(responseBodyReader)
	}
}

func (c *Client) WithResponseBodyReader(responseBodyReader func(io.Reader) error) *Client {
	c.responseBodyReader = responseBodyReader
	return c
}

func (clientOption) WithSafety() funcClientOption {
	return func(c *Client) {
		c.WithSafety()
	}
}

func (c *Client) WithSafety() *Client {
	c.safety = true
	return c
}
