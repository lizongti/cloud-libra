package scheduler

import (
	"fmt"
)

type Scheduler struct {
	opts schedulerOptions

	pipelines []*pipeline
	taskChan  chan *Task
	errorChan chan interface{}
	dieChan   chan struct{}
	exitChan  chan struct{}
}

type pipeline struct {
	dieChan  chan struct{}
	exitChan chan struct{}
}

var defaultScheduler = NewScheduler()
var emptyScheduler *Scheduler

func init() {
	defaultScheduler.WithSafety().WithBackground().Serve()
}

func Default() *Scheduler {
	return defaultScheduler
}

func Empty() *Scheduler {
	return emptyScheduler
}

func NewScheduler(opt ...ApplySchedulerOption) *Scheduler {
	opts := defaultSchedulerOptions

	for _, o := range opt {
		o.apply(&opts)
	}

	s := &Scheduler{
		opts:      opts,
		errorChan: make(chan interface{}, 1),
		dieChan:   make(chan struct{}),
		exitChan:  make(chan struct{}),
	}

	return s
}

func (s *Scheduler) Serve() error {
	s.taskChan = make(chan *Task, s.opts.taskBacklog)

	if s.opts.background {
		go s.serve()
		return nil
	}
	return s.serve()
}

func (s *Scheduler) Close() error {
	for _, p := range s.pipelines {
		s.close(p)
	}
	close(s.dieChan)
	<-s.exitChan
	return nil
}

func (s *Scheduler) serve() (err error) {
	if s.opts.safety {
		defer func() {
			if v := recover(); v != nil {
				err = fmt.Errorf("%v", v)
				if s.opts.errorChan != nil {
					s.opts.errorChan <- err
				}
			}
		}()
	}

	defer close(s.exitChan)

	s.increaseParallel(s.opts.parallel)

	for {
		select {
		case n := <-s.opts.parallelChan:
			s.increaseParallel(n)
		case p := <-s.errorChan:
			panic(p)
		case <-s.dieChan:
			return
		}
	}
}

func (s *Scheduler) increaseParallel(n int) {
	for n > 0 {
		p := &pipeline{
			dieChan:  make(chan struct{}),
			exitChan: make(chan struct{}),
		}
		go func(p *pipeline) {
			s.digest(p)
		}(p)
		s.pipelines = append(s.pipelines, p)
		n--
	}
}

func (s *Scheduler) digest(p *pipeline) {
	defer func() {
		if v := recover(); v != nil {
			err := fmt.Errorf("%v", v)
			s.errorChan <- err
		}
	}()

	defer close(p.exitChan)

	for {
		select {
		case t := <-s.taskChan:
			t.execute()
		case <-p.dieChan:
			return
		}
	}
}

func (s *Scheduler) schedule(t *Task) {
	if s == emptyScheduler {
		t.execute()
		return
	}
	s.taskChan <- t
}

func (s *Scheduler) close(p *pipeline) {
	close(p.dieChan)
	<-p.exitChan
}

func (s *Scheduler) report(r *Report) {
	if s != emptyScheduler && s.opts.reportChan != nil {
		s.opts.reportChan <- r
	}
}

type schedulerOptions struct {
	safety       bool
	background   bool
	errorChan    chan<- error
	taskBacklog  int
	parallel     int
	parallelChan <-chan int
	reportChan   chan<- *Report
}

var defaultSchedulerOptions = schedulerOptions{
	safety:       false,
	background:   false,
	errorChan:    nil,
	taskBacklog:  0,
	parallel:     1,
	parallelChan: nil,
	reportChan:   nil,
}

type ApplySchedulerOption interface {
	apply(*schedulerOptions)
}

type funcSchedulerOption func(*schedulerOptions)

func (f funcSchedulerOption) apply(opt *schedulerOptions) {
	f(opt)
}

type schedulerOption int

var SchedulerOption schedulerOption

func (schedulerOption) Safety() funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.safety = true
	}
}

func (s *Scheduler) WithSafety() *Scheduler {
	SchedulerOption.Safety().apply(&s.opts)
	return s
}

func (schedulerOption) Background() funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.background = true
	}
}

func (s *Scheduler) WithBackground() *Scheduler {
	SchedulerOption.Background().apply(&s.opts)
	return s
}

func (schedulerOption) ErrorChan(errorChan chan<- error) funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.errorChan = errorChan
	}
}

func (s *Scheduler) WithErrorChan(errorChan chan<- error) *Scheduler {
	SchedulerOption.ErrorChan(errorChan).apply(&s.opts)
	return s
}

func (schedulerOption) TaskBacklog(taskBacklog int) funcSchedulerOption {
	return func(s *schedulerOptions) {
		if taskBacklog > 0 {
			s.taskBacklog = taskBacklog
		}
	}
}

func (s *Scheduler) WithTaskBacklog(backlog int) *Scheduler {
	SchedulerOption.TaskBacklog(backlog).apply(&s.opts)
	return s
}

func (schedulerOption) Parallel(parallel int) funcSchedulerOption {
	return func(s *schedulerOptions) {
		if parallel > 0 {
			s.parallel = parallel
		}
	}
}

func (s *Scheduler) WithParallel(parallel int) *Scheduler {
	SchedulerOption.Parallel(parallel).apply(&s.opts)
	return s
}

func (schedulerOption) ReportChan(reportChan chan<- *Report) funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.reportChan = reportChan
	}
}

func (s *Scheduler) WithReportChan(reportChan chan<- *Report) *Scheduler {
	SchedulerOption.ReportChan(reportChan).apply(&s.opts)
	return s
}

func (schedulerOption) ParallelChan(parallelChan <-chan int) funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.parallelChan = parallelChan
	}
}

func (s *Scheduler) WithParallelChan(parallelChan <-chan int) *Scheduler {
	SchedulerOption.ParallelChan(parallelChan).apply(&s.opts)
	return s
}
