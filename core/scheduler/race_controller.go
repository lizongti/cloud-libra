package scheduler

import (
	"fmt"
	"time"
)

type RaceController struct {
	opts      raceControllerOptions
	scheduler *Scheduler
	taskMap   map[string]*Task
	done      int
	failed    int
	timeout   int
}

func NewRaceController(opt ...ApplyRaceControllerOption) *RaceController {
	opts := defaultRaceControllerOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	return &RaceController{
		opts:    opts,
		taskMap: make(map[string]*Task),
	}
}

func (c *RaceController) Serve() error {
	if c.opts.background {
		go c.serve()
		return nil
	}
	return c.serve()
}

func (c *RaceController) serve() (err error) {
	if len(c.opts.tasks) == 0 {
		return nil
	}

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

	taskLength := len(c.opts.tasks)
	reportChan := make(chan *Report, taskLength)

	opt := []ApplySchedulerOption{
		WithSchedulerBackground(),
		WithSchedulerErrorChan(c.opts.errorChan),
		WithSchedulerTaskBacklog(taskLength),
		WithSchedulerParallel(taskLength),
		WithSchedulerReportChan(reportChan),
	}
	if c.opts.safety {
		opt = append(opt, WithSchedulerSafety())
	}
	c.scheduler = NewScheduler(opt...)

	if err := c.scheduler.Serve(); err != nil {
		return err
	}
	defer c.scheduler.Close()

	go func() {
		defer func() {
			if v := recover(); v != nil {
				err := fmt.Errorf("%v", v)
				if c.opts.errorChan != nil {
					c.opts.errorChan <- err
				}
			}
		}()
		for _, task := range c.opts.tasks {
			task.Publish(c.scheduler)
		}
	}()

	for _, task := range c.opts.tasks {
		c.taskMap[task.ID()] = task
	}

	defer func() {
		c.timeout = taskLength - c.done - c.failed
	}()

	if c.opts.timeout > 0 {
		timer := time.After(c.opts.timeout)
		for {
			select {
			case r := <-reportChan:
				switch {
				case r.State == TaskStateDone:
					c.done++
					if c.opts.doneFunc != nil {
						c.opts.doneFunc(c.taskMap[r.ID])
					}
				case r.State == TaskStateFailed:
					c.failed++
					if c.opts.failedFunc != nil {
						c.opts.failedFunc(c.taskMap[r.ID])
					}
				}
				if c.done+c.failed == taskLength {
					return
				}
			case <-timer:
				return
			}
		}
	} else {
		for {
			select {
			case r := <-reportChan:
				switch {
				case r.State == TaskStateDone:
					c.done++
					if c.opts.doneFunc != nil {
						c.opts.doneFunc(c.taskMap[r.ID])
					}
				case r.State == TaskStateFailed:
					c.failed++
					if c.opts.failedFunc != nil {
						c.opts.failedFunc(c.taskMap[r.ID])
					}
				}
				if c.done+c.failed == taskLength {
					return
				}
			}
		}
	}

}

func (c *RaceController) Size() int {
	return len(c.opts.tasks)
}

func (c *RaceController) Done() int {
	return c.done
}

func (c *RaceController) Failed() int {
	return c.failed
}

func (c *RaceController) Timeout() int {
	return c.timeout
}

type raceControllerOptions struct {
	safety     bool
	background bool
	errorChan  chan<- error
	timeout    time.Duration
	tasks      []*Task
	doneFunc   func(*Task)
	failedFunc func(*Task)
}

var defaultRaceControllerOptions = raceControllerOptions{
	safety:     false,
	background: false,
	errorChan:  nil,
	timeout:    0,
	tasks:      nil,
	doneFunc:   nil,
	failedFunc: nil,
}

type ApplyRaceControllerOption interface {
	apply(*raceControllerOptions)
}

type funcRaceControllerOption func(*raceControllerOptions)

func (f funcRaceControllerOption) apply(opt *raceControllerOptions) {
	f(opt)
}

func WithRaceSafety() funcRaceControllerOption {
	return func(c *raceControllerOptions) {
		c.safety = true
	}
}

func WithRaceBackground() funcRaceControllerOption {
	return func(c *raceControllerOptions) {
		c.background = true
	}
}

func WithRaceErrorChan(errorChan chan<- error) funcRaceControllerOption {
	return func(c *raceControllerOptions) {
		c.errorChan = errorChan
	}
}

func WithRaceTasks(tasks ...*Task) funcRaceControllerOption {
	return func(c *raceControllerOptions) {
		c.tasks = tasks
	}
}

func WithRaceTimeout(timeout time.Duration) funcRaceControllerOption {
	return func(c *raceControllerOptions) {
		c.timeout = timeout
	}
}

func WithRaceDoneFunc(doneFunc func(*Task)) funcRaceControllerOption {
	return func(c *raceControllerOptions) {
		c.doneFunc = doneFunc
	}
}

func WithRaceFailedFunc(failedFunc func(*Task)) funcRaceControllerOption {
	return func(c *raceControllerOptions) {
		c.failedFunc = failedFunc
	}
}
