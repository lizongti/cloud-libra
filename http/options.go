package http

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strings"
)

type options struct {
	ContentType        string
	Timeout            int
	Retry              int
	Proxy              string
	Form               url.Values
	Body               io.Reader
	ResponseBodyReader func(io.Reader) error
}

func newOptions() *options {
	return &options{
		Form: make(url.Values),
	}
}

func (c *Client) WithTimeout(timeout int) *Client {
	c.Timeout = timeout
	return c
}

func (c *Client) WithRetry(retry int) *Client {
	c.Retry = retry
	return c
}

func (c *Client) WithProxy(proxy string) *Client {
	c.Proxy = proxy
	return c
}

func (c *Client) WithContentType(contentType string) *Client {
	c.ContentType = contentType
	return c
}

func (c *Client) WithForm(form url.Values) *Client {
	c.Form = form
	return c
}

func (c *Client) WithParam(key string, value string) *Client {
	c.Form.Set(key, value)
	return c
}

func (c *Client) WithBody(body interface{}) *Client {
	switch body := body.(type) {
	case string:
		c.Body = strings.NewReader(body)
	case []byte:
		c.Body = bytes.NewReader(body)
	case io.Reader:
		c.Body = body
	default:
		c.Body = strings.NewReader(fmt.Sprint("%v", body))
	}
	return c
}

func (c *Client) WithResponseBodyReader(responseBodyReader func(io.Reader) error) *Client {
	c.ResponseBodyReader = responseBodyReader
	return c
}
