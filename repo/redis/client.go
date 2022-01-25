package redis

import (
	"context"
	"time"

	"github.com/aceaura/libra/core/cast"
	"github.com/gomodule/redigo/redis"
)

type Client struct {
	opts clientOptions
}

func NewClient(opt ...ApplyClientOption) *Client {
	opts := defaultClientOptions

	for _, o := range opt {
		o.apply(&opts)
	}

	c := &Client{
		opts: opts,
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
	return redis.DoContext(conn, c.opts.context, commandName, args...)
}

func (c *Client) Command(commands ...interface{}) ([]string, error) {
	reply, err := c.Do(commands[0].(string), commands[1:]...)
	if err != nil {
		return nil, err
	}

	return parseReply(reply), err
}

func (c *Client) dial() (redis.Conn, error) {
	conn, err := redis.DialContext(c.opts.context, "tcp", c.opts.addr,
		redis.DialPassword(c.opts.password),
		redis.DialDatabase(c.opts.db),
		redis.DialConnectTimeout(c.opts.connectTimeout),
		redis.DialReadTimeout(c.opts.readTimeout),
		redis.DialWriteTimeout(c.opts.writeTimeout),
	)
	if err != nil {
		return nil, err
	}
	if _, err := redis.DoContext(conn, c.opts.context, "SELECT", c.opts.db); err != nil {
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
	return redis.DoContext(conn, p.client.opts.context, commandName, args...)
}

func (p *Pool) Command(commands ...interface{}) ([]string, error) {
	reply, err := p.Do(commands[0].(string), commands[1:]...)
	if err != nil {
		return nil, err
	}

	return parseReply(reply), err
}

type clientOptions struct {
	context  context.Context
	addr     string
	password string
	db       int

	maxActive int
	maxIdle   int

	keepAlive      time.Duration
	connectTimeout time.Duration
	readTimeout    time.Duration
	writeTimeout   time.Duration
	idleTimeout    time.Duration
}

var defaultClientOptions = clientOptions{
	context:  context.Background(),
	addr:     "localhost:6379",
	password: "",
	db:       0,

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

func (fro funcClientOption) apply(so *clientOptions) {
	fro(so)
}

type clientOption int

var ClientOption clientOption

func (clientOption) Context(context context.Context) funcClientOption {
	return func(c *clientOptions) {
		c.context = context
	}
}

func (c *Client) WithContext(context context.Context) *Client {
	ClientOption.Context(context).apply(&c.opts)
	return c
}

func (clientOption) Addr(addr string) funcClientOption {
	return func(c *clientOptions) {
		c.addr = addr
	}
}

func (c *Client) WithAddr(addr string) *Client {
	ClientOption.Addr(addr).apply(&c.opts)
	return c
}

func (clientOption) Password(password string) funcClientOption {
	return func(c *clientOptions) {
		c.password = password
	}
}

func (c *Client) WithPassword(password string) *Client {
	ClientOption.Password(password).apply(&c.opts)
	return c
}

func (clientOption) DB(db int) funcClientOption {
	return func(c *clientOptions) {
		c.db = db
	}
}

func (c *Client) WithDB(db int) *Client {
	ClientOption.DB(db).apply(&c.opts)
	return c
}

func (clientOption) MaxActive(maxActive int) funcClientOption {
	return func(c *clientOptions) {
		c.maxActive = maxActive
	}
}

func (c *Client) WithMaxActive(maxActive int) *Client {
	ClientOption.MaxActive(maxActive).apply(&c.opts)
	return c
}

func (clientOption) MaxIdle(maxIdle int) funcClientOption {
	return func(c *clientOptions) {
		c.maxIdle = maxIdle
	}
}

func (c *Client) WithMaxIdle(maxIdle int) *Client {
	ClientOption.MaxIdle(maxIdle).apply(&c.opts)
	return c
}

func (clientOption) ConnectTimeout(connectTimeout time.Duration) funcClientOption {
	return func(c *clientOptions) {
		c.connectTimeout = connectTimeout
	}
}

func (c *Client) WithConnectTimeout(connectTimeout time.Duration) *Client {
	ClientOption.ConnectTimeout(connectTimeout).apply(&c.opts)
	return c
}

func (clientOption) ReadTimeout(readTimeout time.Duration) funcClientOption {
	return func(c *clientOptions) {
		c.readTimeout = readTimeout
	}
}

func (c *Client) WithReadTimeout(readTimeout time.Duration) *Client {
	ClientOption.ReadTimeout(readTimeout).apply(&c.opts)
	return c
}

func (clientOption) WriteTimeout(writeTimeout time.Duration) funcClientOption {
	return func(c *clientOptions) {
		c.writeTimeout = writeTimeout
	}
}

func (c *Client) WithWriteTimeout(writeTimeout time.Duration) *Client {
	ClientOption.WriteTimeout(writeTimeout).apply(&c.opts)
	return c
}

func (clientOption) IdleTimeout(idleTimeout time.Duration) funcClientOption {
	return func(c *clientOptions) {
		c.idleTimeout = idleTimeout
	}
}

func (c *Client) WithIdleTimeout(idleTimeout time.Duration) *Client {
	ClientOption.IdleTimeout(idleTimeout).apply(&c.opts)
	return c
}
