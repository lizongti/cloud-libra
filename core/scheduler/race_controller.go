package scheduler

import (
	"fmt"
	"time"
)

type RaceController struct {
	opts    raceControllerOptions
	done    int
	failed  int
	timeout int
}

func NewRaceController(opt ...ApplyRaceControllerOption) *RaceController {
	opts := defaultRaceControllerOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	return &RaceController{
		opts: opts,
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
	if len(c.opts.taskMap) == 0 {
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

	taskLength := len(c.opts.taskMap)
	errorChan := make(chan error, taskLength)
	reportChan := make(chan *Report, taskLength)

	s := NewScheduler()
	if c.opts.safety {
		s.WithSafety()
	}
	s.WithBackground()
	s.WithErrorChan(errorChan)
	s.WithTaskBacklog(taskLength)
	s.WithParallel(taskLength)
	s.WithReportChan(reportChan)
	if err := s.Serve(); err != nil {
		return err
	}
	defer s.Close()

	go func() {
		for _, task := range c.opts.taskMap {
			task.Publish(s)
		}
	}()

	defer func() {
		c.timeout = taskLength - c.done - c.failed
	}()

	if c.opts.timeout > 0 {
		timer := time.After(c.opts.timeout)
		for {
			select {
			case e := <-errorChan:
				c.opts.errorFunc(e)
			case r := <-reportChan:
				switch {
				case r.State == TaskStateDone:
					c.done++
					c.opts.doneFunc(c.opts.taskMap[r.ID])
				case r.State == TaskStateFailed:
					c.failed++
					c.opts.failedFunc(c.opts.taskMap[r.ID])
				case c.done+c.failed == taskLength:
					return
				}
			case <-timer:
				return
			}
		}
	} else {
		for {
			select {
			case e := <-errorChan:
				c.opts.errorFunc(e)
			case r := <-reportChan:
				switch {
				case r.State == TaskStateDone:
					c.done++
					c.opts.doneFunc(c.opts.taskMap[r.ID])
				case r.State == TaskStateFailed:
					c.failed++
					c.opts.failedFunc(c.opts.taskMap[r.ID])
				}
				if c.done+c.failed == taskLength {
					return
				}
			}
		}
	}

}

func (c *RaceController) Size() int {
	return len(c.opts.taskMap)
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
	taskMap    map[string]*Task
	errorFunc  func(error)
	doneFunc   func(*Task)
	failedFunc func(*Task)
}

var defaultRaceControllerOptions = raceControllerOptions{
	safety:     false,
	background: false,
	errorChan:  nil,
	timeout:    0,
	taskMap:    map[string]*Task{},
	errorFunc:  func(error) {},
	doneFunc:   func(*Task) {},
	failedFunc: func(*Task) {},
}

type ApplyRaceControllerOption interface {
	apply(*raceControllerOptions)
}

type funcRaceControllerOption func(*raceControllerOptions)

func (f funcRaceControllerOption) apply(opt *raceControllerOptions) {
	f(opt)
}

type raceControllerOption int

var RaceControllerOption raceControllerOption

func (raceControllerOption) Safety() funcRaceControllerOption {
	return func(c *raceControllerOptions) {
		c.safety = true
	}
}

func (c *RaceController) WithSafety() *RaceController {
	RaceControllerOption.Safety().apply(&c.opts)
	return c
}

func (raceControllerOption) Background() funcRaceControllerOption {
	return func(c *raceControllerOptions) {
		c.background = true
	}
}

func (c *RaceController) WithBackground() *RaceController {
	RaceControllerOption.Background().apply(&c.opts)
	return c
}

func (raceControllerOption) ErrorChan(errorChan chan<- error) funcRaceControllerOption {
	return func(c *raceControllerOptions) {
		c.errorChan = errorChan
	}
}

func (c *RaceController) WithErrorChan(errorChan chan<- error) *RaceController {
	RaceControllerOption.ErrorChan(errorChan).apply(&c.opts)
	return c
}

func (raceControllerOption) Task(tasks ...*Task) funcRaceControllerOption {
	return func(c *raceControllerOptions) {
		for _, task := range tasks {
			c.taskMap[task.ID()] = task
		}
	}
}

func (c *RaceController) WithTask(tasks ...*Task) *RaceController {
	RaceControllerOption.Task(tasks...).apply(&c.opts)
	return c
}

func (raceControllerOption) Timeout(timeout time.Duration) funcRaceControllerOption {
	return func(c *raceControllerOptions) {
		c.timeout = timeout
	}
}

func (c *RaceController) WithTimeout(timeout time.Duration) *RaceController {
	RaceControllerOption.Timeout(timeout).apply(&c.opts)
	return c
}

func (raceControllerOption) ErrorFunc(errorFunc func(error)) funcRaceControllerOption {
	return func(c *raceControllerOptions) {
		c.errorFunc = errorFunc
	}
}

func (c *RaceController) WithErrorFunc(errorFunc func(error)) *RaceController {
	RaceControllerOption.ErrorFunc(errorFunc).apply(&c.opts)
	return c
}

func (raceControllerOption) DoneFunc(doneFunc func(*Task)) funcRaceControllerOption {
	return func(c *raceControllerOptions) {
		c.doneFunc = doneFunc
	}
}

func (c *RaceController) WithDoneFunc(doneFunc func(*Task)) *RaceController {
	RaceControllerOption.DoneFunc(doneFunc).apply(&c.opts)
	return c
}

func (raceControllerOption) FailedFunc(failedFunc func(*Task)) funcRaceControllerOption {
	return func(c *raceControllerOptions) {
		c.failedFunc = failedFunc
	}
}

func (c *RaceController) WithFailedFunc(failedFunc func(*Task)) *RaceController {
	RaceControllerOption.FailedFunc(failedFunc).apply(&c.opts)
	return c
}
