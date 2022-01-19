package coroutine

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCoroutineIsDead       = errors.New("coroutine is dead")
	ErrCoroutineNotSuspended = errors.New("coroutine is not suspended")
	ErrCoroutineTimeout      = errors.New("coroutine timeout")
)

type CoroutineStateType int

const (
	// CoroutineStateCreated means coroutine is created and not started.
	CoroutineStateCreated = iota

	// CoroutineStateSuspended means coroutine is started and yielded.
	CoroutineStateSuspended

	// CoroutineStateRunning means coroutine is started and running.
	CoroutineStateRunning

	// CoroutineStateDead means coroutine not created or ended.
	CoroutineStateDead
)

var coroutineStateName = map[CoroutineStateType]string{
	CoroutineStateCreated:   "created",
	CoroutineStateSuspended: "suspended",
	CoroutineStateRunning:   "running",
	CoroutineStateDead:      "dead",
}

func (t CoroutineStateType) String() string {
	if s, ok := coroutineStateName[t]; ok {
		return s
	}
	return fmt.Sprintf("coroutineStateName=%d?", int(t))
}

type (
	// ID is the unique identifier for coroutine
	ID = string

	Coroutines = sync.Map
)

var (
	coroutines Coroutines
)

// Coroutine is a simulator struct for coroutine
type Coroutine struct {
	opts        coroutineOptions
	id          ID
	status      CoroutineStateType
	inCh        chan []interface{}
	outCh       chan []interface{}
	mutexStatus *sync.Mutex
	mutexResume *sync.Mutex
}

func NewCoroutine(opt ...ApplyCoroutineOption) *Coroutine {
	opts := defaultCoroutineOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	c := &Coroutine{
		opts:        opts,
		id:          uuid.NewString(),
		status:      CoroutineStateCreated,
		inCh:        make(chan []interface{}, 1),
		outCh:       make(chan []interface{}, 1),
		mutexStatus: &sync.Mutex{},
		mutexResume: &sync.Mutex{},
	}

	return c
}

// Start wraps and starts a coroutine up.
// It is thread-safe, and it should be called before other funcs.
func Start(f func(c *Coroutine) error) error {
	return Wrap(func(c *Coroutine, args ...interface{}) error {
		return f(c)
	}).Call()
}

// Create wraps and yields a coroutine with no args, waits for a resume.
// It is not thread-safe, and it should be called before other funcs.
// Call `Resume` after `Create` to start up a coroutine.
func Create(f func(c *Coroutine, inData ...interface{}) error) (*Coroutine, chan error) {
	c := Wrap(func(c *Coroutine, args ...interface{}) error {
		outData, err := c.Yield()
		if err != nil {
			return err
		}
		return f(c, outData...)
	})
	errChan := make(chan error, 1)
	go func() {
		errChan <- c.Call()
	}()
	return c, errChan
}

// Wrap wraps a coroutine and waits for a startup.
// It is thread-safe, and it should be called before other funcs.
// Call `Call` after `Wrap` to start up a coroutine.
func Wrap(f func(c *Coroutine, args ...interface{}) error) *Coroutine {
	c := NewCoroutine().WithFunc(f)
	coroutines.Store(c.ID(), c)
	return c
}

// Call launch a coroutine that is already wrapped.
// It is not thread-safe, and it can only be called beside after Wrap.
// Call `Call` After `Wrap` to start up a coroutine.
func Call(id ID, args ...interface{}) error {
	v, ok := coroutines.Load(id)
	if !ok {
		return ErrCoroutineIsDead
	}
	c := v.(*Coroutine)
	return c.Call(args...)
}

// Resume continues a suspend ID, passing data in and out.
// It is thread-safe, and it can only be called in other Goroutine.
// Call `Resume` after `Create` to start up a coroutine.
// Call `Resume` after `Yield` to continue a coroutine.
func Resume(id ID, inData ...interface{}) ([]interface{}, error) {
	v, ok := coroutines.Load(id)
	if !ok {
		return nil, ErrCoroutineIsDead
	}
	c := v.(*Coroutine)
	return c.Resume(inData...)
}

// TryResume likes Resume, but checks status instead of waiting for status.
// It is thread-safe, and it can only be called in other Goroutine.
// Call `TryResume` after `Create` to start up a coroutine.
// Call `TryResume` after `Yield` to continue a coroutine.
func TryResume(id ID, inData ...interface{}) ([]interface{}, error) {
	v, ok := coroutines.Load(id)
	if !ok {
		return nil, ErrCoroutineIsDead
	}
	c := v.(*Coroutine)
	return c.TryResume(inData...)
}

// Yield suspends a running coroutine, passing data in and out.
// It is not thread-safe, and it can only be called in entity.fn.
// Call `Resume`, `TryResume` or `AsyncResume`
// after `Yield` to continue a coroutine.
func Yield(id ID, outData ...interface{}) ([]interface{}, error) {
	v, ok := coroutines.Load(id)
	if !ok {
		return nil, ErrCoroutineIsDead
	}
	c := v.(*Coroutine)
	return c.Yield(outData...)
}

// Status shows the status of a coroutine.
// It is thread-safe, and it can be called in any Goroutine.
// Call `Status` anywhere you need.
func Status(id ID) CoroutineStateType {
	v, ok := coroutines.Load(id)
	if !ok {
		return CoroutineStateDead
	}
	c := v.(*Coroutine)
	return c.Status()
}

func (c *Coroutine) ID() ID {
	return c.id
}

func (c *Coroutine) Call(args ...interface{}) error {
	c.writeSyncStatus(CoroutineStateRunning)

	return func() (err error) {
		defer func() {
			if v := recover(); v != nil {
				err = fmt.Errorf("coroutine %s error:%v", c.id, v)
			}
		}()
		defer func() {
			coroutines.Delete(c.id)
		}()

		return c.opts.f(c, args...)
	}()
}

func (c *Coroutine) Resume(inData ...interface{}) ([]interface{}, error) {
	c.mutexResume.Lock()
	defer c.mutexResume.Unlock()
	if c.readSyncStatus() == CoroutineStateDead {
		return nil, ErrCoroutineIsDead
	}
	return c.resume(inData)
}

func (c *Coroutine) TryResume(inData ...interface{}) ([]interface{}, error) {
	c.mutexResume.Lock()
	defer c.mutexResume.Unlock()
	if c.readSyncStatus() != CoroutineStateSuspended {
		return nil, ErrCoroutineNotSuspended
	}
	return c.resume(inData)
}

func (c *Coroutine) Yield(outData ...interface{}) ([]interface{}, error) {
	c.writeSyncStatus(CoroutineStateSuspended)
	inData, err := c.yield(outData)
	if err != nil {
		return nil, err
	}
	c.writeSyncStatus(CoroutineStateRunning)
	return inData, nil
}

func (c *Coroutine) Status() CoroutineStateType {
	return c.readSyncStatus()
}

func (c *Coroutine) writeSyncStatus(status CoroutineStateType) {
	c.mutexStatus.Lock()
	defer c.mutexStatus.Unlock()
	c.status = status
}

func (c *Coroutine) readSyncStatus() CoroutineStateType {
	c.mutexStatus.Lock()
	defer c.mutexStatus.Unlock()
	return c.status
}

func (c *Coroutine) resume(inData []interface{}) ([]interface{}, error) {
	var outData []interface{}

	select {
	case outData = <-c.outCh:
		break
	case <-time.After(c.opts.timeout):
		return nil, ErrCoroutineTimeout
	}

	select {
	case c.inCh <- inData:
		break
	case <-time.After(c.opts.timeout):
		return nil, ErrCoroutineTimeout
	}

	return outData, nil
}

func (c *Coroutine) yield(outData []interface{}) ([]interface{}, error) {
	var inData []interface{}

	select {
	case c.outCh <- outData:
		break
	case <-time.After(c.opts.timeout):
		c.writeSyncStatus(CoroutineStateDead)
		return nil, ErrCoroutineTimeout
	}

	select {
	case inData = <-c.inCh:
		break
	case <-time.After(c.opts.timeout):
		c.writeSyncStatus(CoroutineStateDead)
		return nil, ErrCoroutineTimeout
	}

	return inData, nil
}

type coroutineOptions struct {
	timeout time.Duration
	f       func(*Coroutine, ...interface{}) error
}

var defaultCoroutineOptions = coroutineOptions{
	timeout: 30 * time.Second,
	f:       func(*Coroutine, ...interface{}) error { return nil },
}

type ApplyCoroutineOption interface {
	apply(*coroutineOptions)
}

type funcCoroutineOption func(*coroutineOptions)

func (fco funcCoroutineOption) apply(id *coroutineOptions) {
	fco(id)
}

type coroutineOption int

var CoroutineOption coroutineOption

func (coroutineOption) Timeout(timeout time.Duration) funcCoroutineOption {
	return func(c *coroutineOptions) {
		c.timeout = timeout
	}
}

func (c *Coroutine) WithTimeout(timeout time.Duration) *Coroutine {
	CoroutineOption.Timeout(timeout).apply(&c.opts)
	return c
}

func (coroutineOption) Func(f func(*Coroutine, ...interface{}) error) funcCoroutineOption {
	return func(c *coroutineOptions) {
		c.f = f
	}
}

func (c *Coroutine) WithFunc(f func(*Coroutine, ...interface{}) error) *Coroutine {
	CoroutineOption.Func(f).apply(&c.opts)
	return c
}
