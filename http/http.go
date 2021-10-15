package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	GET  = "GET"
	POST = "POST"
	HEAD = "HEAD"
)

type Client struct {
	*options
	client http.Client
}

func NewClient() *Client {
	return &Client{
		options: newOptions(),
	}
}

func Get(url string, params ...interface{}) (*http.Response, []byte, error) {
	if len(params) == 0 {
		return NewClient().Get(url)
	}
	return params[0].(*Client).Get(url)
}

func (c *Client) Get(url string) (*http.Response, []byte, error) {
	return c.Do(GET, url)
}

func Post(url string, params ...interface{}) (*http.Response, []byte, error) {
	if len(params) == 0 {
		return NewClient().Post(url)
	}
	return params[0].(*Client).Post(url)
}

func (c *Client) Post(url string) (*http.Response, []byte, error) {
	return c.Do(POST, url)
}

func Head(url string, params ...interface{}) (*http.Response, []byte, error) {
	if len(params) == 0 {
		return NewClient().Head(url)
	}
	return params[0].(*Client).Head(url)
}

func (c *Client) Head(url string) (*http.Response, []byte, error) {
	return c.Do(HEAD, url)
}

func (c *Client) Do(method string, url string) (*http.Response, []byte, error) {
	c.init()
	return c.requestWithRetry(func() (*http.Response, error) {
		if len(c.Form) > 0 {
			url = fmt.Sprintf("%s?%s", url, c.Form.Encode())
		}
		req, err := http.NewRequest(method, url, c.Body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", c.ContentType)
		return c.client.Do(req)
	})
}

func (c *Client) requestWithRetry(fn func() (*http.Response, error)) (resp *http.Response, body []byte, err error) {
	for times := 0; times < c.Retry; times++ {
		resp, body, err = c.request(fn)
		if err != nil {
			return nil, nil, err
		}
		if resp.StatusCode != 200 {
			continue
		}
	}
	return resp, body, err
}

func (c *Client) request(fn func() (*http.Response, error)) (resp *http.Response, body []byte, err error) {
	resp, err = fn()
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return nil, nil, err
	}
	if c.ResponseBodyReader != nil {
		err = c.ResponseBodyReader(resp.Body)
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
	if c.Proxy != "" {
		fixedURL, err := url.Parse(c.Proxy)
		if err != nil {
			return err
		}
		c.client.Transport = &http.Transport{Proxy: http.ProxyURL(fixedURL)}
	}
	if c.Timeout > 0 {
		c.client.Timeout = time.Duration(c.Timeout)
	}
	if c.Retry < 1 {
		c.Retry = 1
	}
	if c.ContentType == "" {
		c.ContentType = "text/plain"
	}
	return nil
}
