package scheduler

import (
	"fmt"
	"math"
	"time"
)

type Controller struct {
	opts     controllerOptions
	dieChan  chan struct{}
	exitChan chan struct{}
}

func NewController(opt ...ApplyControllerOption) *Controller {
	opts := defaultControllerOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	return &Controller{
		opts:     opts,
		dieChan:  make(chan struct{}),
		exitChan: make(chan struct{}),
	}
}

func (c *Controller) Serve(reportChan <-chan *Report, parallelChan chan<- int) error {
	if c.opts.background {
		go c.serve(reportChan, parallelChan)
		return nil
	}
	return c.serve(reportChan, parallelChan)
}

func (c *Controller) serve(reportChan <-chan *Report, parallelChan chan<- int) (err error) {
	if c.opts.safety {
		defer func() {
			if v := recover(); v != nil {
				err = fmt.Errorf("%v", v)
				if c.opts.errorChan != nil {
					c.opts.errorChan <- err
				}
			}
		}()
	}

	defer close(c.exitChan)

	var (
		tickerChan = time.NewTicker(c.opts.parallelTick).C
		stateMap   = make(map[TaskStateType]int)
		finished   int
		tpsMax     int
	)

	for {
		select {
		case r := <-reportChan:
			if r.Progress == 0 {
				stateMap[r.State]++
			}
			if r.State == TaskStateDone || r.State == TaskStateFailed {
				finished++
			}

		case <-tickerChan:
			tps := int(float64(finished) * float64(time.Second) / float64(c.opts.parallelTick))
			tpsMax = int(math.Max(float64(tps), float64(tpsMax)))
			finished = 0
			if stateMap[TaskStatePending]-stateMap[TaskStateRunning] > 0 {
				if tpsMax < c.opts.tpsLimit || c.opts.tpsLimit < 0 {
					parallelChan <- c.opts.parallelIncrease
				}
			}
		case <-c.dieChan:
			return
		}
	}
}

func (c *Controller) Close() {
	close(c.dieChan)
	<-c.exitChan
}

type controllerOptions struct {
	safety           bool
	background       bool
	errorChan        chan<- error
	parallelTick     time.Duration
	parallelIncrease int
	tpsLimit         int
}

var defaultControllerOptions = controllerOptions{
	safety:           false,
	background:       false,
	errorChan:        nil,
	parallelTick:     time.Second,
	parallelIncrease: 1,
	tpsLimit:         -1,
}

type ApplyControllerOption interface {
	apply(*controllerOptions)
}

type funcControllerOption func(*controllerOptions)

func (f funcControllerOption) apply(opt *controllerOptions) {
	f(opt)
}

type controllerOption int

var ControllerOption controllerOption

func (controllerOption) Safety() funcControllerOption {
	return func(c *controllerOptions) {
		c.safety = true
	}
}

func (c *Controller) Safety() *Controller {
	ControllerOption.Safety().apply(&c.opts)
	return c
}

func (controllerOption) Background() funcControllerOption {
	return func(c *controllerOptions) {
		c.background = true
	}
}

func (c *Controller) WithBackground() *Controller {
	ControllerOption.Background().apply(&c.opts)
	return c
}

func (controllerOption) ErrorChan(errorChan chan<- error) funcControllerOption {
	return func(c *controllerOptions) {
		c.errorChan = errorChan
	}
}

func (c *Controller) ErrorChan(errorChan chan<- error) *Controller {
	ControllerOption.ErrorChan(errorChan).apply(&c.opts)
	return c
}

func (controllerOption) ParallelTick(parallelTick time.Duration) funcControllerOption {
	return func(c *controllerOptions) {
		c.parallelTick = parallelTick
	}
}

func (c *Controller) WithParallelTick(parallelTick time.Duration) *Controller {
	ControllerOption.ParallelTick(parallelTick).apply(&c.opts)
	return c
}

func (controllerOption) ParallelIncrease(parallelIncrease int) funcControllerOption {
	return func(c *controllerOptions) {
		if parallelIncrease > 0 {
			c.parallelIncrease = parallelIncrease
		}
	}
}

func (c *Controller) WithParallelIncrease(parallelIncrease int) *Controller {
	ControllerOption.ParallelIncrease(parallelIncrease).apply(&c.opts)
	return c
}

func (controllerOption) TPSLimit(tpsLimit int) funcControllerOption {
	return func(c *controllerOptions) {
		if tpsLimit > 0 {
			c.tpsLimit = tpsLimit
		}
	}
}

func (c *Controller) WithTPSLimit(tpsLimit int) *Controller {
	ControllerOption.TPSLimit(tpsLimit).apply(&c.opts)
	return c
}
