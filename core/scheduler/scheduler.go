package scheduler

import (
	"fmt"
)

type Scheduler struct {
	opts schedulerOptions

	pipelines []*pipeline
	taskChan  chan *Task
	panicChan chan interface{}
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
	defaultScheduler.Safety().Background().Serve()
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
		panicChan: make(chan interface{}, 1),
		dieChan:   make(chan struct{}),
		exitChan:  make(chan struct{}),
	}

	return s
}

func (s *Scheduler) Serve() error {
	s.taskChan = make(chan *Task, s.opts.backlog)

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
	if s.opts.errorFunc != nil {
		defer func() {
			s.opts.errorFunc(err)
			err = nil
		}()
	}

	if s.opts.safety {
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("%v", e)
			}
		}()
	}

	defer close(s.exitChan)

	s.increaseParallel(s.opts.parallel)

	for {
		select {
		case n := <-s.opts.parallelChan:
			s.increaseParallel(n)
		case p := <-s.panicChan:
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
		if p := recover(); p != nil {
			s.panicChan <- p
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
	backlog      int
	parallel     int
	parallelChan <-chan int
	background   bool
	safety       bool
	errorFunc    func(error)
	reportChan   chan<- *Report
}

var defaultSchedulerOptions = schedulerOptions{}

type ApplySchedulerOption interface {
	apply(*schedulerOptions)
}

type funcSchedulerOption func(*schedulerOptions)

func (fso funcSchedulerOption) apply(so *schedulerOptions) {
	fso(so)
}

type schedulerOption int

var SchedulerOption schedulerOption

func (schedulerOption) Backlog(backlog int) funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.backlog = backlog
	}
}

func (s *Scheduler) Backlog(backlog int) *Scheduler {
	SchedulerOption.Backlog(backlog).apply(&s.opts)
	return s
}

func (schedulerOption) Parallel(parallel int) funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.parallel = parallel
	}
}

func (s *Scheduler) Parallel(parallel int) *Scheduler {
	SchedulerOption.Parallel(parallel).apply(&s.opts)
	return s
}

func (schedulerOption) Background() funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.background = true
	}
}

func (s *Scheduler) Background() *Scheduler {
	SchedulerOption.Background().apply(&s.opts)
	return s
}

func (schedulerOption) Safety() funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.safety = true
	}
}

func (s *Scheduler) Safety() *Scheduler {
	SchedulerOption.Safety().apply(&s.opts)
	return s
}

func (schedulerOption) ErrorFunc(errorFunc func(error)) funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.errorFunc = errorFunc
	}
}

func (s *Scheduler) ErrorFunc(errorFunc func(error)) *Scheduler {
	SchedulerOption.ErrorFunc(errorFunc).apply(&s.opts)
	return s
}

func (schedulerOption) ReportChan(reportChan chan<- *Report) funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.reportChan = reportChan
	}
}

func (s *Scheduler) ReportChan(reportChan chan<- *Report) *Scheduler {
	SchedulerOption.ReportChan(reportChan).apply(&s.opts)
	return s
}

func (schedulerOption) ParallelChan(parallelChan <-chan int) funcSchedulerOption {
	return func(s *schedulerOptions) {
		s.parallelChan = parallelChan
	}
}

func (s *Scheduler) ParallelChan(parallelChan <-chan int) *Scheduler {
	SchedulerOption.ParallelChan(parallelChan).apply(&s.opts)
	return s
}
