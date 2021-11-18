package scheduler

import (
	"fmt"
	"math"
	"time"
)

type TPSController struct {
	opts     tpsControllerOptions
	dieChan  chan struct{}
	exitChan chan struct{}
}

func NewTPSController(opt ...ApplyTPSControllerOption) *TPSController {
	opts := defaultTPSControllerOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	return &TPSController{
		opts:     opts,
		dieChan:  make(chan struct{}),
		exitChan: make(chan struct{}),
	}
}

func (c *TPSController) Serve(reportChan <-chan *Report, parallelChan chan<- int) error {
	if c.opts.background {
		go c.serve(reportChan, parallelChan)
		return nil
	}
	return c.serve(reportChan, parallelChan)
}

func (c *TPSController) serve(reportChan <-chan *Report, parallelChan chan<- int) (err error) {
	if c.opts.safety {
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("%v", e)
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

func (c *TPSController) Close() {
	close(c.dieChan)
	<-c.exitChan
}

type tpsControllerOptions struct {
	safety           bool
	background       bool
	errorChan        chan<- error
	parallelTick     time.Duration
	parallelIncrease int
	tpsLimit         int
}

var defaultTPSControllerOptions = tpsControllerOptions{
	safety:           false,
	background:       false,
	errorChan:        nil,
	parallelTick:     time.Second,
	parallelIncrease: 1,
	tpsLimit:         -1,
}

type ApplyTPSControllerOption interface {
	apply(*tpsControllerOptions)
}

type funcTPSControllerOption func(*tpsControllerOptions)

func (fco funcTPSControllerOption) apply(co *tpsControllerOptions) {
	fco(co)
}

type tpsControllerOption int

var ControllerOption tpsControllerOption

func (tpsControllerOption) Safety() funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		c.safety = true
	}
}

func (c *TPSController) Safety() *TPSController {
	ControllerOption.Safety().apply(&c.opts)
	return c
}

func (tpsControllerOption) Background() funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		c.background = true
	}
}

func (c *TPSController) WithBackground() *TPSController {
	ControllerOption.Background().apply(&c.opts)
	return c
}

func (tpsControllerOption) ErrorChan(errorChan chan<- error) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		c.errorChan = errorChan
	}
}

func (c *TPSController) ErrorChan(errorChan chan<- error) *TPSController {
	ControllerOption.ErrorChan(errorChan).apply(&c.opts)
	return c
}

func (tpsControllerOption) ParallelTick(parallelTick time.Duration) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		c.parallelTick = parallelTick
	}
}

func (c *TPSController) WithParallelTick(parallelTick time.Duration) *TPSController {
	ControllerOption.ParallelTick(parallelTick).apply(&c.opts)
	return c
}

func (tpsControllerOption) ParallelIncrease(parallelIncrease int) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		if parallelIncrease > 0 {
			c.parallelIncrease = parallelIncrease
		}
	}
}

func (c *TPSController) WithParallelIncrease(parallelIncrease int) *TPSController {
	ControllerOption.ParallelIncrease(parallelIncrease).apply(&c.opts)
	return c
}

func (tpsControllerOption) TPSLimit(tpsLimit int) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		if tpsLimit > 0 {
			c.tpsLimit = tpsLimit
		}
	}
}

func (c *TPSController) WithTPSLimit(tpsLimit int) *TPSController {
	ControllerOption.TPSLimit(tpsLimit).apply(&c.opts)
	return c
}
