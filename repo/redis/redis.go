package redis

import (
	"github.com/go-redis/redis/v8"
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
	options := &Options{
		Addr:         r.opts.addr,
		Password:     r.opts.password,
		DB:           r.opts.db,
		PoolSize:     r.opts.size,
		MinIdleConns: r.opts.idle,
	}
	return redis.NewClient(options)
}

type redisOptions struct {
	addr     string
	password string
	db       int
	retry    int
	size     int
	idle     int
}

var defaultRedisOptions = redisOptions{
	addr:     "localhost:6379",
	password: "",
	db:       0,
	retry:    0,
	size:     0,
	idle:     0,
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

func (redisOption) Retry(retry int) funcRedisOption {
	return func(r *redisOptions) {
		r.retry = retry
	}
}

func (r *Redis) WithRetry(retry int) *Redis {
	RedisOption.Retry(retry).apply(&r.opts)
	return r
}

func (redisOption) Size(size int) funcRedisOption {
	return func(r *redisOptions) {
		r.size = size
	}
}

func (r *Redis) WithSize(size int) *Redis {
	RedisOption.Size(size).apply(&r.opts)
	return r
}

func (redisOption) Idle(idle int) funcRedisOption {
	return func(r *redisOptions) {
		r.idle = idle
	}
}

func (r *Redis) WithIdle(idle int) *Redis {
	RedisOption.Idle(idle).apply(&r.opts)
	return r
}
