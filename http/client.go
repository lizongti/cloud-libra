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

type Client struct {
	*clientOpt
	client http.Client
}

func NewClient(options ...clientOption) *Client {
	c := &Client{
		clientOpt: newClientOpt(options),
	}
	c.doOpt(c)
	return c
}

func Get(url string, options ...clientOption) (*http.Response, []byte, error) {
	return NewClient(options...).Get(url)
}

func (c *Client) Get(url string) (*http.Response, []byte, error) {
	return c.Do(GET, url)
}

func Post(url string, options ...clientOption) (*http.Response, []byte, error) {
	return NewClient(options...).Post(url)
}

func (c *Client) Post(url string) (*http.Response, []byte, error) {
	return c.Do(POST, url)
}

func Head(url string, options ...clientOption) (*http.Response, []byte, error) {
	return NewClient(options...).Head(url)
}

func (c *Client) Head(url string) (*http.Response, []byte, error) {
	return c.Do(HEAD, url)
}

func Do(method string, url string, options ...clientOption) (*http.Response, []byte, error) {
	return NewClient(options...).Do(method, url)
}

func (c *Client) Do(method string, url string) (resp *http.Response, body []byte, err error) {
	c.init()
	if c.safety {
		defer func() {
			if v := recover(); v != nil {
				err = v.(error)
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

func (c *Client) init() error {
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

type clientOption func(*Client)
type clientOptions []clientOption

type clientOpt struct {
	clientOptions
	protocol           string
	contentType        string
	timeout            time.Duration
	retry              int
	proxy              string
	form               url.Values
	body               io.Reader
	responseBodyReader func(io.Reader) error
	safety             bool
}

func newClientOpt(options []clientOption) *clientOpt {
	return &clientOpt{
		clientOptions: options,
		form:          make(url.Values),
	}
}

func (opt *clientOpt) doOpt(c *Client) {
	for _, option := range opt.clientOptions {
		option(c)
	}
}

func WithProtocol(protocol string) clientOption {
	return func(c *Client) {
		c.WithProtocol(protocol)
	}
}

func (c *Client) WithProtocol(protocol string) *Client {
	c.protocol = protocol
	return c
}

func WithTimeout(timeout time.Duration) clientOption {
	return func(c *Client) {
		c.WithTimeout(timeout)
	}
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.timeout = timeout
	return c
}

func WithRetry(retry int) clientOption {
	return func(c *Client) {
		c.WithRetry(retry)
	}
}

func (c *Client) WithRetry(retry int) *Client {
	c.retry = retry
	return c
}

func WithProxy(proxy string) clientOption {
	return func(c *Client) {
		c.WithProxy(proxy)
	}
}

func (c *Client) WithProxy(proxy string) *Client {
	c.proxy = proxy
	return c
}

func WithContentType(contentType string) clientOption {
	return func(c *Client) {
		c.WithContentType(contentType)
	}
}

func (c *Client) WithContentType(contentType string) *Client {
	c.contentType = contentType
	return c
}

func WithForm(form url.Values) clientOption {
	return func(c *Client) {
		c.WithForm(form)
	}
}

func (c *Client) WithForm(form url.Values) *Client {
	c.form = form
	return c
}

func WithParam(key string, value string) clientOption {
	return func(c *Client) {
		c.WithParam(key, value)
	}
}

func (c *Client) WithParam(key string, value string) *Client {
	c.form.Set(key, value)
	return c
}

func WithBody(body interface{}) clientOption {
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

func WithResponseBodyReader(responseBodyReader func(io.Reader) error) clientOption {
	return func(c *Client) {
		c.WithResponseBodyReader(responseBodyReader)
	}
}

func (c *Client) WithResponseBodyReader(responseBodyReader func(io.Reader) error) *Client {
	c.responseBodyReader = responseBodyReader
	return c
}

func WithClientSafety() clientOption {
	return func(c *Client) {
		c.WithClientSafety()
	}
}

func (c *Client) WithClientSafety() *Client {
	c.safety = true
	return c
}
