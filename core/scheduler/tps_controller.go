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
		opts:      opts,
		scheduler: NewScheduler(),
		dieChan:   make(chan struct{}),
		exitChan:  make(chan struct{}),
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

	if c.opts.safety {
		c.scheduler.WithSafety()
	}
	c.scheduler.WithBackground()
	c.scheduler.WithErrorChan(c.opts.errorChan)
	c.scheduler.WithTaskBacklog(c.opts.taskBacklog)
	c.scheduler.WithParallel(c.opts.parallel)
	c.scheduler.WithReportChan(reportChan)
	c.scheduler.WithParallelChan(parallelChan)

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

type tpsControllerOption int

var TPSControllerOption tpsControllerOption

func (tpsControllerOption) Safety() funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		c.safety = true
	}
}

func (c *TPSController) Safety() *TPSController {
	TPSControllerOption.Safety().apply(&c.opts)
	return c
}

func (tpsControllerOption) Background() funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		c.background = true
	}
}

func (c *TPSController) WithBackground() *TPSController {
	TPSControllerOption.Background().apply(&c.opts)
	return c
}

func (tpsControllerOption) ErrorChan(errorChan chan<- error) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		c.errorChan = errorChan
	}
}

func (c *TPSController) WithErrorChan(errorChan chan<- error) *TPSController {
	TPSControllerOption.ErrorChan(errorChan).apply(&c.opts)
	return c
}

func (tpsControllerOption) Parallel(parallel int) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		c.parallel = parallel
	}
}

func (c *TPSController) WithParallel(parallel int) *TPSController {
	TPSControllerOption.Parallel(parallel).apply(&c.opts)
	return c
}

func (tpsControllerOption) ParallelTick(parallelTick time.Duration) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		c.parallelTick = parallelTick
	}
}

func (c *TPSController) WithParallelTick(parallelTick time.Duration) *TPSController {
	TPSControllerOption.ParallelTick(parallelTick).apply(&c.opts)
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
	TPSControllerOption.ParallelIncrease(parallelIncrease).apply(&c.opts)
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
	TPSControllerOption.TPSLimit(tpsLimit).apply(&c.opts)
	return c
}

func (tpsControllerOption) TaskBacklog(taskBacklog int) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		if taskBacklog > 0 {
			c.taskBacklog = taskBacklog
		}
	}
}

func (c *TPSController) WithTaskBacklog(taskBacklog int) *TPSController {
	TPSControllerOption.TaskBacklog(taskBacklog).apply(&c.opts)
	return c
}

func (tpsControllerOption) ParallelBacklog(parallelBacklog int) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		if parallelBacklog > 0 {
			c.parallelBacklog = parallelBacklog
		}
	}
}

func (c *TPSController) WithParallelBacklog(parallelBacklog int) *TPSController {
	TPSControllerOption.ParallelBacklog(parallelBacklog).apply(&c.opts)
	return c
}

func (tpsControllerOption) ReportBacklog(reportBacklog int) funcTPSControllerOption {
	return func(c *tpsControllerOptions) {
		if reportBacklog > 0 {
			c.reportBacklog = reportBacklog
		}
	}
}

func (c *TPSController) WithReportBacklog(reportBacklog int) *TPSController {
	TPSControllerOption.ReportBacklog(reportBacklog).apply(&c.opts)
	return c
}
