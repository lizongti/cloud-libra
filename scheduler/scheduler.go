package scheduler

import (
	"fmt"
)

type Scheduler struct {
	opts         []schedulerOpt
	backlog      int
	parallel     int
	parallelChan <-chan int
	background   bool
	safety       bool
	errorFunc    func(error)
	reportChan   chan<- *Report
	pipelines    []*pipeline
	taskChan     chan *Task
	panicChan    chan interface{}
	dieChan      chan struct{}
	exitChan     chan struct{}
}

type pipeline struct {
	dieChan  chan struct{}
	exitChan chan struct{}
}

func NewScheduler(opts ...schedulerOpt) *Scheduler {
	return &Scheduler{opts: opts}
}

func Default() *Scheduler {
	return NewScheduler()
}

func (s *Scheduler) Serve() error {
	s.init()
	if s.background {
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

func (s *Scheduler) init() {
	s.panicChan = make(chan interface{}, 1)
	s.dieChan = make(chan struct{})
	s.exitChan = make(chan struct{})

	for _, opt := range s.opts {
		opt(s)
	}

	s.taskChan = make(chan *Task, s.backlog)
	if s.parallel == 0 {
		s.parallel = 1
	}
}

func (s *Scheduler) serve() (err error) {
	if s.errorFunc != nil {
		defer func() {
			s.errorFunc(err)
			err = nil
		}()
	}

	if s.safety {
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("%v", e)
			}
		}()
	}

	defer close(s.exitChan)

	s.increaseParallel(s.parallel)

	for {
		select {
		case n := <-s.parallelChan:
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
	s.taskChan <- t
}

func (s *Scheduler) close(p *pipeline) {
	close(p.dieChan)
	<-p.exitChan
}

func (s *Scheduler) needReport() bool {
	return s.reportChan != nil
}

func (s *Scheduler) report(r *Report) {
	s.reportChan <- r
}

type schedulerOpt func(*Scheduler)
type schedulerOption struct{}

var SchedulerOption schedulerOption

func (schedulerOption) WithBacklog(backlog int) schedulerOpt {
	return func(s *Scheduler) {
		s.backlog = backlog
	}
}

func (s *Scheduler) WithBacklog(backlog int) *Scheduler {
	s.opts = append(s.opts, SchedulerOption.WithBacklog(backlog))
	return s
}

func (schedulerOption) WithParallel(parallel int) schedulerOpt {
	return func(s *Scheduler) {
		s.parallel = parallel
	}
}

func (s *Scheduler) WithParallel(parallel int) *Scheduler {
	s.opts = append(s.opts, SchedulerOption.WithParallel(parallel))
	return s
}

func (schedulerOption) WithBackground() schedulerOpt {
	return func(s *Scheduler) {
		s.background = true
	}
}

func (s *Scheduler) WithBackground() *Scheduler {
	s.opts = append(s.opts, SchedulerOption.WithBackground())
	return s
}

func (schedulerOption) WithSchdedulerSafety() schedulerOpt {
	return func(s *Scheduler) {
		s.safety = true
	}
}

func (s *Scheduler) WithSchdedulerSafety() *Scheduler {
	s.opts = append(s.opts, SchedulerOption.WithSchdedulerSafety())
	return s
}

func (schedulerOption) WithErrorFunc(errorFunc func(error)) schedulerOpt {
	return func(s *Scheduler) {
		s.errorFunc = errorFunc
	}
}

func (s *Scheduler) WithErrorFunc(errorFunc func(error)) *Scheduler {
	s.opts = append(s.opts, SchedulerOption.WithErrorFunc(errorFunc))
	return s
}

func (schedulerOption) WithReportChan(reportChan chan<- *Report) schedulerOpt {
	return func(s *Scheduler) {
		s.reportChan = reportChan
	}
}

func (s *Scheduler) WithReportChan(reportChan chan<- *Report) *Scheduler {
	s.opts = append(s.opts, SchedulerOption.WithReportChan(reportChan))
	return s
}

func (schedulerOption) WithParallelChan(parallelChan <-chan int) schedulerOpt {
	return func(s *Scheduler) {
		s.parallelChan = parallelChan
	}
}

func (s *Scheduler) WithParallelChan(parallelChan <-chan int) *Scheduler {
	s.opts = append(s.opts, SchedulerOption.WithParallelChan(parallelChan))
	return s
}
