package coroutine

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type CoroutineStateType int

const (
	// CoroutineStateCreated means ID is created and not started.
	CoroutineStateCreated = iota

	// CoroutineStateSuspended means ID is started and yielded.
	CoroutineStateSuspended

	// CoroutineStateRunning means ID is started and running.
	CoroutineStateRunning

	// CoroutineStateDead means ID not created or ended.
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

// Start wraps and starts a ID up.
// It is thread-safe, and it should be called before other funcs.
func Start(f func(c *Coroutine) error) error {
	return Call(Wrap(func(c *Coroutine, args ...interface{}) error {
		return f(c)
	}).ID())
}

// Create wraps and yields a ID with no args, waits for a resume.
// It is not thread-safe, and it should be called before other funcs.
// Call `Resume` after `Create` to start up a ID.
func Create(f func(c *Coroutine, inData ...interface{}) error) *Coroutine {
	c := Wrap(func(c *Coroutine, args ...interface{}) error {
		return f(c, c.Yield()...)
	})
	go func() {
		if err := c.Call(); err != nil {
			panic(err)
		}
	}()
	return c
}

// Wrap wraps a ID and waits for a startup.
// It is thread-safe, and it should be called before other funcs.
// Call `Call` after `Wrap` to start up a ID.
func Wrap(f func(c *Coroutine, args ...interface{}) error) *Coroutine {
	c := NewCoroutine().WithFunc(f)
	coroutines.Store(c.ID(), c)
	return c
}

// Call launch a ID that is already wrapped.
// It is not thread-safe, and it can only be called beside after Wrap.
// Call `Call` After `Wrap` to start up a ID.
func Call(id ID, args ...interface{}) error {
	v, ok := coroutines.Load(id)
	if !ok {
		panic(fmt.Errorf("ID %s is not found", id))
	}
	c := v.(*Coroutine)
	return c.Call(args...)
}

// Resume continues a suspend ID, passing data in and out.
// It is thread-safe, and it can only be called in other Goroutine.
// Call `Resume` after `Create` to start up a ID.
// Call `Resume` after `Yield` to continue a ID.
func Resume(id ID, inData ...interface{}) ([]interface{}, bool) {
	v, ok := coroutines.Load(id)
	if !ok {
		panic(fmt.Errorf("coroutine %s is missing or dead", id))
	}
	c := v.(*Coroutine)
	return c.Resume(inData...)
}

// TryResume likes Resume, but checks status instead of waiting for status.
// It is thread-safe, and it can only be called in other Goroutine.
// Call `TryResume` after `Create` to start up a ID.
// Call `TryResume` after `Yield` to continue a ID.
func TryResume(id ID, inData ...interface{}) ([]interface{}, bool) {
	v, ok := coroutines.Load(id)
	if !ok {
		panic(fmt.Errorf("coroutine %s is missing or dead", id))
	}
	c := v.(*Coroutine)
	return c.TryResume(inData...)
}

// Yield suspends a running coroutine, passing data in and out.
// It is not thread-safe, and it can only be called in entity.fn.
// Call `Resume`, `TryResume` or `AsyncResume`
// after `Yield` to continue a ID.
func Yield(id ID, outData ...interface{}) []interface{} {
	v, ok := coroutines.Load(id)
	if !ok {
		panic(fmt.Errorf("coroutine %s is missing or dead", id))
	}
	c := v.(*Coroutine)
	return c.Yield(outData...)
}

// Status shows the status of a ID.
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

func (c *Coroutine) Resume(inData ...interface{}) ([]interface{}, bool) {
	c.mutexResume.Lock()
	defer c.mutexResume.Unlock()
	if c.readSyncStatus() == CoroutineStateDead {
		return nil, false
	}
	outData := c.resume(inData)
	return outData, true
}

func (c *Coroutine) TryResume(inData ...interface{}) ([]interface{}, bool) {
	c.mutexResume.Lock()
	defer c.mutexResume.Unlock()
	if c.readSyncStatus() != CoroutineStateSuspended {
		return nil, false
	}
	outData := c.resume(inData)
	return outData, true
}

func (c *Coroutine) Yield(outData ...interface{}) []interface{} {
	c.writeSyncStatus(CoroutineStateSuspended)
	inData := c.yield(outData)
	c.writeSyncStatus(CoroutineStateRunning)
	return inData
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

func (c *Coroutine) resume(inData []interface{}) []interface{} {
	var outData []interface{}

	select {
	case outData = <-c.outCh:
		break
	case <-time.After(c.opts.timeout):
		panic(fmt.Errorf("ID %s suspended for too long", c.id))
	}

	select {
	case c.inCh <- inData:
		break
	case <-time.After(c.opts.timeout):
		panic(fmt.Errorf("ID %s suspended for too long", c.id))
	}

	return outData
}

func (c *Coroutine) yield(outData []interface{}) []interface{} {
	var inData []interface{}

	select {
	case c.outCh <- outData:
		break
	case <-time.After(c.opts.timeout):
		c.writeSyncStatus(CoroutineStateDead)
		panic(fmt.Errorf("ID %s suspended for too long", c.id))
	}

	select {
	case inData = <-c.inCh:
		break
	case <-time.After(c.opts.timeout):
		c.writeSyncStatus(CoroutineStateDead)
		panic(fmt.Errorf("ID %s suspended for too long", c.id))
	}

	return inData
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
