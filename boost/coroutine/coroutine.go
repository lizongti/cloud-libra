package coroutine

import (
	"context"
	"errors"
	"fmt"
	"sync"

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
	f           func(*Coroutine, ...interface{}) error
	ctx         context.Context
	id          ID
	state       CoroutineStateType
	inChan      chan []interface{}
	outChan     chan []interface{}
	mutexStatus *sync.Mutex
	mutexResume *sync.Mutex
}

func NewCoroutine(ctx context.Context, f func(*Coroutine, ...interface{}) error) *Coroutine {
	c := &Coroutine{
		f:           f,
		ctx:         ctx,
		id:          uuid.NewString(),
		state:       CoroutineStateCreated,
		inChan:      make(chan []interface{}, 1),
		outChan:     make(chan []interface{}, 1),
		mutexStatus: &sync.Mutex{},
		mutexResume: &sync.Mutex{},
	}

	return c
}

// Start wraps and starts a coroutine up.
// It is thread-safe, and it should be called before other funcs.
func Start(ctx context.Context, f func(c *Coroutine) error) error {
	return Wrap(ctx, func(c *Coroutine, args ...interface{}) error {
		return f(c)
	}).Call()
}

// Create wraps and yields a coroutine with no args, waits for a resume.
// It is not thread-safe, and it should be called before other funcs.
// Call `Resume` after `Create` to start up a coroutine.
func Create(ctx context.Context, f func(*Coroutine, ...interface{}) error) (*Coroutine, chan error) {
	c := Wrap(ctx, func(c *Coroutine, args ...interface{}) error {
		out, err := c.Yield()
		if err != nil {
			return err
		}
		return f(c, out...)
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
func Wrap(ctx context.Context, f func(c *Coroutine, args ...interface{}) error) *Coroutine {
	c := NewCoroutine(ctx, f)
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
func Resume(id ID, in ...interface{}) ([]interface{}, error) {
	v, ok := coroutines.Load(id)
	if !ok {
		return nil, ErrCoroutineIsDead
	}
	c := v.(*Coroutine)
	return c.Resume(in...)
}

// TryResume likes Resume, but checks status instead of waiting for status.
// It is thread-safe, and it can only be called in other Goroutine.
// Call `TryResume` after `Create` to start up a coroutine.
// Call `TryResume` after `Yield` to continue a coroutine.
func TryResume(id ID, in ...interface{}) ([]interface{}, error) {
	v, ok := coroutines.Load(id)
	if !ok {
		return nil, ErrCoroutineIsDead
	}
	c := v.(*Coroutine)
	return c.TryResume(in...)
}

// Yield suspends a running coroutine, passing data in and out.
// It is not thread-safe, and it can only be called in entity.fn.
// Call `Resume`, `TryResume` or `AsyncResume`
// after `Yield` to continue a coroutine.
func Yield(id ID, out ...interface{}) ([]interface{}, error) {
	v, ok := coroutines.Load(id)
	if !ok {
		return nil, ErrCoroutineIsDead
	}
	c := v.(*Coroutine)
	return c.Yield(out...)
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
	defer func() {
		coroutines.Delete(c.id)
	}()

	c.setStatus(CoroutineStateRunning)

	errChan := make(chan error)
	defer close(errChan)

	go func() {
		defer func() {
			if v := recover(); v != nil {
				err := fmt.Errorf("coroutine %s error:%v", c.id, v)
				errChan <- err
			}
		}()

		errChan <- c.f(c, args...)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	case <-c.ctx.Done():
		c.setStatus(CoroutineStateDead)
		return c.ctx.Err()
	}
	return nil
}

func (c *Coroutine) Resume(in ...interface{}) ([]interface{}, error) {
	c.mutexResume.Lock()
	defer c.mutexResume.Unlock()
	if c.status() == CoroutineStateDead {
		return nil, ErrCoroutineIsDead
	}
	return c.resume(in)
}

func (c *Coroutine) TryResume(in ...interface{}) ([]interface{}, error) {
	c.mutexResume.Lock()
	defer c.mutexResume.Unlock()
	if c.status() != CoroutineStateSuspended {
		return nil, ErrCoroutineNotSuspended
	}
	return c.resume(in)
}

func (c *Coroutine) Yield(out ...interface{}) ([]interface{}, error) {
	c.setStatus(CoroutineStateSuspended)
	defer c.setStatus(CoroutineStateRunning)
	if c.status() == CoroutineStateDead {
		return nil, ErrCoroutineIsDead
	}

	return c.yield(out)
}

func (c *Coroutine) Status() CoroutineStateType {
	return c.status()
}

func (c *Coroutine) setStatus(state CoroutineStateType) {
	c.mutexStatus.Lock()
	defer c.mutexStatus.Unlock()
	c.state = state
}

func (c *Coroutine) status() CoroutineStateType {
	c.mutexStatus.Lock()
	defer c.mutexStatus.Unlock()
	return c.state
}

func (c *Coroutine) resume(in []interface{}) ([]interface{}, error) {
	errChan := make(chan error)
	defer close(errChan)

	var out []interface{}

	go func() {
		defer func() {
			if v := recover(); v != nil {
				err := fmt.Errorf("coroutine %s error:%v", c.id, v)
				errChan <- err
			}
		}()

		out = <-c.outChan
		c.inChan <- in
		errChan <- nil
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	case <-c.ctx.Done():
		c.setStatus(CoroutineStateDead)
		return nil, c.ctx.Err()
	}

	return out, nil
}

func (c *Coroutine) yield(out []interface{}) ([]interface{}, error) {
	c.outChan <- out
	return <-c.inChan, nil
}
