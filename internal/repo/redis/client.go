package redis

import (
	"context"
	"net/url"
	"time"

	"github.com/cloudlibraries/libra/internal/boost/cast"
	"github.com/gomodule/redigo/redis"
)

type Client struct {
	opts clientOptions
	url  string
}

func NewClient(url string, opt ...ApplyClientOption) *Client {
	opts := defaultClientOptions

	for _, o := range opt {
		o.apply(&opts)
	}

	c := &Client{
		opts: opts,
		url:  url,
	}

	return c
}

func (c *Client) Pool() *Pool {
	return &Pool{
		client: c,
		pool: &redis.Pool{
			MaxIdle:     c.opts.maxIdle,
			MaxActive:   c.opts.maxActive,
			IdleTimeout: c.opts.idleTimeout,
			Dial:        c.dial,
		},
	}
}

func (c *Client) Do(commandName string, args ...interface{}) (interface{}, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return redis.DoContext(conn, c.opts.ctx, commandName, args...)
}

func (c *Client) Command(commands ...interface{}) ([]string, error) {
	reply, err := c.Do(commands[0].(string), commands[1:]...)
	if err != nil {
		return nil, err
	}

	return parseReply(reply), err
}

func (c *Client) dial() (redis.Conn, error) {
	u, err := url.Parse(c.url)
	if err != nil {
		return nil, err
	}

	addr := u.Host
	db := cast.ToInt(u.Path[1:])
	password, _ := u.User.Password()

	conn, err := redis.DialContext(c.opts.ctx, "tcp", addr,
		redis.DialPassword(password),
		redis.DialDatabase(db),
		redis.DialConnectTimeout(c.opts.connectTimeout),
		redis.DialReadTimeout(c.opts.readTimeout),
		redis.DialWriteTimeout(c.opts.writeTimeout),
	)
	if err != nil {
		return nil, err
	}
	if _, err := redis.DoContext(conn, c.opts.ctx, "SELECT", db); err != nil {
		return nil, err
	}
	return conn, nil
}

func parseReply(reply interface{}) []string {
	result := make([]string, 0)
	switch reply := reply.(type) {
	case []interface{}:
		result = cast.ToStringSlice(reply)
	default:
		result = append(result, cast.ToString(reply))
	}
	return result
}

type Pool struct {
	client *Client
	pool   *redis.Pool
}

func (p *Pool) Do(commandName string, args ...interface{}) (interface{}, error) {
	conn := p.pool.Get()
	defer conn.Close()
	return redis.DoContext(conn, p.client.opts.ctx, commandName, args...)
}

func (p *Pool) Command(commands ...interface{}) ([]string, error) {
	reply, err := p.Do(commands[0].(string), commands[1:]...)
	if err != nil {
		return nil, err
	}

	return parseReply(reply), err
}

type clientOptions struct {
	ctx context.Context

	maxActive int
	maxIdle   int

	keepAlive      time.Duration
	connectTimeout time.Duration
	readTimeout    time.Duration
	writeTimeout   time.Duration
	idleTimeout    time.Duration
}

var defaultClientOptions = clientOptions{
	ctx: context.Background(),

	maxActive: 0,
	maxIdle:   0,

	connectTimeout: 30 * time.Second,
	readTimeout:    30 * time.Second,
	writeTimeout:   30 * time.Second,
	idleTimeout:    60 * time.Second,
}

type ApplyClientOption interface {
	apply(*clientOptions)
}

type funcClientOption func(*clientOptions)

func (f funcClientOption) apply(opt *clientOptions) {
	f(opt)
}

func WithClientContext(ctx context.Context) funcClientOption {
	return func(c *clientOptions) {
		c.ctx = ctx
	}
}

func WithMaxActive(maxActive int) funcClientOption {
	return func(c *clientOptions) {
		c.maxActive = maxActive
	}
}

func WithClientMaxIdle(maxIdle int) funcClientOption {
	return func(c *clientOptions) {
		c.maxIdle = maxIdle
	}
}

func WithClientConnectTimeout(connectTimeout time.Duration) funcClientOption {
	return func(c *clientOptions) {
		c.connectTimeout = connectTimeout
	}
}

func WithClientReadTimeout(readTimeout time.Duration) funcClientOption {
	return func(c *clientOptions) {
		c.readTimeout = readTimeout
	}
}

func WithClientWriteTimeout(writeTimeout time.Duration) funcClientOption {
	return func(c *clientOptions) {
		c.writeTimeout = writeTimeout
	}
}

func WithClientIdleTimeout(idleTimeout time.Duration) funcClientOption {
	return func(c *clientOptions) {
		c.idleTimeout = idleTimeout
	}
}
