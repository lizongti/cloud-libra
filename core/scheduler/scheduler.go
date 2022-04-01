package scheduler

import (
	"fmt"
)

type Scheduler struct {
	opts schedulerOptions

	pipelines []*pipeline
	taskChan  chan *Task
	dieChan   chan struct{}
	exitChan  chan struct{}
}

type pipeline struct {
	dieChan  chan struct{}
	exitChan chan struct{}
}

var defaultScheduler = NewScheduler(
	WithSchedulerSafety(),
	WithSchedulerBackground(),
)
var emptyScheduler *Scheduler

func init() {
	defaultScheduler.Serve()
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
		opts:     opts,
		dieChan:  make(chan struct{}),
		exitChan: make(chan struct{}),
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
			if s.opts.errorChan != nil {
				s.opts.errorChan <- err
			}
		}
	}()

	defer close(p.exitChan)

	for {
		select {
		case t := <-s.taskChan:
			s.execute(t)
		case <-p.dieChan:
			return
		}
	}
}

func (s *Scheduler) schedule(t *Task) {
	if s == emptyScheduler {
		s.execute(t)
		return
	}
	s.taskChan <- t
}

func (s *Scheduler) execute(t *Task) {
	if err := t.Execute(); err != nil {
		if s.opts.errorChan != nil {
			s.opts.errorChan <- err
		}
	}
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

func WithSchedulerSafety() funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.safety = true
	}
}

func WithSchedulerBackground() funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.background = true
	}
}

func WithSchedulerErrorChan(errorChan chan<- error) funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.errorChan = errorChan
	}
}

func WithSchedulerTaskBacklog(taskBacklog int) funcSchedulerOption {
	return func(s *schedulerOptions) {
		if taskBacklog > 0 {
			s.taskBacklog = taskBacklog
		}
	}
}

func WithSchedulerParallel(parallel int) funcSchedulerOption {
	return func(s *schedulerOptions) {
		if parallel > 0 {
			s.parallel = parallel
		}
	}
}

func WithSchedulerReportChan(reportChan chan<- *Report) funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.reportChan = reportChan
	}
}

func WithSchedulerParallelChan(parallelChan <-chan int) funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.parallelChan = parallelChan
	}
}
