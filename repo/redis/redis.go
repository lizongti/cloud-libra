package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/cast"
)

type Redis struct {
	opts redisOptions
}

func NewRedis(opt ...ApplyRedisOption) *Redis {
	opts := defaultRedisOptions

	for _, o := range opt {
		o.apply(&opts)
	}

	r := &Redis{
		opts: opts,
	}

	return r
}

func (r *Redis) Client() *Client {
	return &Client{
		Pool: &redis.Pool{
			MaxIdle:     r.opts.maxIdle,
			MaxActive:   r.opts.maxActive,
			IdleTimeout: r.opts.idleTimeout,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", r.opts.addr, redis.DialPassword(r.opts.password),
					redis.DialDatabase(r.opts.db),
					redis.DialConnectTimeout(r.opts.connectTimeout),
					redis.DialReadTimeout(r.opts.readTimeout),
					redis.DialWriteTimeout(r.opts.writeTimeout),
				)
				if err != nil {
					return nil, err
				}
				if _, err := conn.Do("SELECT", r.opts.db); err != nil {
					return nil, err
				}
				return conn, nil
			},
		},
	}
}

type Client struct {
	*redis.Pool
}

func (c *Client) Do(commandName string, args ...interface{}) ([]string, error) {
	conn := c.Get()
	defer conn.Close()
	reply, err := conn.Do(commandName, args...)
	if err != nil {
		return nil, err
	}

	return c.parseReply(reply), err
}

func (c *Client) parseReply(reply interface{}) []string {
	result := make([]string, 0)
	switch reply := reply.(type) {
	case []interface{}:
		result = cast.ToStringSlice(reply)
	default:
		result = append(result, cast.ToString(reply))
	}
	return result
}

type redisOptions struct {
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

var defaultRedisOptions = redisOptions{
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

type ApplyRedisOption interface {
	apply(*redisOptions)
}

type funcRedisOption func(*redisOptions)

func (fro funcRedisOption) apply(so *redisOptions) {
	fro(so)
}

type redisOption int

var RedisOption redisOption

func (redisOption) Addr(addr string) funcRedisOption {
	return func(r *redisOptions) {
		r.addr = addr
	}
}

func (r *Redis) WithAddr(addr string) *Redis {
	RedisOption.Addr(addr).apply(&r.opts)
	return r
}

func (redisOption) Password(password string) funcRedisOption {
	return func(r *redisOptions) {
		r.password = password
	}
}

func (r *Redis) WithPassword(password string) *Redis {
	RedisOption.Password(password).apply(&r.opts)
	return r
}

func (redisOption) DB(db int) funcRedisOption {
	return func(r *redisOptions) {
		r.db = db
	}
}

func (r *Redis) WithDB(db int) *Redis {
	RedisOption.DB(db).apply(&r.opts)
	return r
}

func (redisOption) MaxActive(maxActive int) funcRedisOption {
	return func(r *redisOptions) {
		r.maxActive = maxActive
	}
}

func (r *Redis) WithMaxActive(maxActive int) *Redis {
	RedisOption.MaxActive(maxActive).apply(&r.opts)
	return r
}

func (redisOption) MaxIdle(maxIdle int) funcRedisOption {
	return func(r *redisOptions) {
		r.maxIdle = maxIdle
	}
}

func (r *Redis) WithMaxIdle(maxIdle int) *Redis {
	RedisOption.MaxIdle(maxIdle).apply(&r.opts)
	return r
}

func (redisOption) ConnectTimeout(connectTimeout time.Duration) funcRedisOption {
	return func(r *redisOptions) {
		r.connectTimeout = connectTimeout
	}
}

func (r *Redis) WithConnectTimeout(connectTimeout time.Duration) *Redis {
	RedisOption.ConnectTimeout(connectTimeout).apply(&r.opts)
	return r
}

func (redisOption) ReadTimeout(readTimeout time.Duration) funcRedisOption {
	return func(r *redisOptions) {
		r.readTimeout = readTimeout
	}
}

func (r *Redis) WithReadTimeout(readTimeout time.Duration) *Redis {
	RedisOption.ReadTimeout(readTimeout).apply(&r.opts)
	return r
}

func (redisOption) WriteTimeout(writeTimeout time.Duration) funcRedisOption {
	return func(r *redisOptions) {
		r.writeTimeout = writeTimeout
	}
}

func (r *Redis) WithWriteTimeout(writeTimeout time.Duration) *Redis {
	RedisOption.WriteTimeout(writeTimeout).apply(&r.opts)
	return r
}

func (redisOption) IdleTimeout(idleTimeout time.Duration) funcRedisOption {
	return func(r *redisOptions) {
		r.idleTimeout = idleTimeout
	}
}

func (r *Redis) WithIdleTimeout(idleTimeout time.Duration) *Redis {
	RedisOption.IdleTimeout(idleTimeout).apply(&r.opts)
	return r
}
