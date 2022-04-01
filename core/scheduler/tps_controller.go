package scheduler

import (
	"fmt"
	"math"
	"time"
)

type TPSController struct {
	opts      tpsControllerOptions
	scheduler *Scheduler
	dieChan   chan struct{}
	exitChan  chan struct{}
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

func (c *TPSController) Serve() error {
	if c.opts.background {
		go c.serve()
		return nil
	}
	return c.serve()
}

func (c *TPSController) Scheduler() *Scheduler {
	return c.scheduler
}

func (c *TPSController) serve() (err error) {
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

	reportChan := make(chan *Report, c.opts.reportBacklog)
	parallelChan := make(chan int, c.opts.parallelBacklog)

	opt := []ApplySchedulerOption{
		WithBackground(),
		WithErrorChan(c.opts.errorChan),
		WithTaskBacklog(c.opts.taskBacklog),
		WithParallel(c.opts.parallel),
		WithReportChan(reportChan),
		WithParallelChan(parallelChan),
	}
	if c.opts.safety {
		opt = append(opt, WithSafety())
	}
	c.scheduler = NewScheduler(opt...)

	if err := c.scheduler.Serve(); err != nil {
		return err
	}
	defer c.scheduler.Close()

	tickerChan := time.NewTicker(c.opts.parallelTick).C
	stateMap := make(map[TaskStateType]int)
	finished := int(0)
	tpsMax := int(0)

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
	parallel         int
	parallelTick     time.Duration
	parallelIncrease int
	tpsLimit         int
	taskBacklog      int
	parallelBacklog  int
	reportBacklog    int
}

var defaultTPSControllerOptions = tpsControllerOptions{
	safety:           false,
	background:       false,
	errorChan:        nil,
	parallel:         1,
	parallelTick:     time.Second,
	parallelIncrease: 1,
	tpsLimit:         -1,
	taskBacklog:      0,
	parallelBacklog:  1,
	reportBacklog:    0,
}

type ApplyTPSControllerOption interface {
	apply(*tpsControllerOptions)
}

type funcTPSControllerOption func(*tpsControllerOptions)

func (f funcTPSControllerOption) apply(opt *tpsControllerOptions) {
	f(opt)
}

func WithTPSSafety() funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		c.safety = true
	}
}

func WithTPSBackground() funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		c.background = true
	}
}

func WithTPSErrorChan(errorChan chan<- error) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		c.errorChan = errorChan
	}
}

func WithTPSParallel(parallel int) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		c.parallel = parallel
	}
}

func WithTPSParallelTick(parallelTick time.Duration) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		c.parallelTick = parallelTick
	}
}

func WithTPSParallelIncrease(parallelIncrease int) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		if parallelIncrease > 0 {
			c.parallelIncrease = parallelIncrease
		}
	}
}

func WithTPSLimit(tpsLimit int) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		if tpsLimit > 0 {
			c.tpsLimit = tpsLimit
		}
	}
}

func WithTPSTaskBacklog(taskBacklog int) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		if taskBacklog > 0 {
			c.taskBacklog = taskBacklog
		}
	}
}

func WithTPSParallelBacklog(parallelBacklog int) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		if parallelBacklog > 0 {
			c.parallelBacklog = parallelBacklog
		}
	}
}

func WithTPSReportBacklog(reportBacklog int) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		if reportBacklog > 0 {
			c.reportBacklog = reportBacklog
		}
	}
}
