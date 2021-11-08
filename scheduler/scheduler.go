package scheduler

import (
	"fmt"
)

type Scheduler struct {
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

var scheduler = NewScheduler()

func init() {
	scheduler.WithSafety().WithBackground().Serve()
}

func Default() *Scheduler {
	return scheduler
}

func NewScheduler(opts ...funcSchedulerOption) *Scheduler {
	s := &Scheduler{
		panicChan: make(chan interface{}, 1),
		dieChan:   make(chan struct{}),
		exitChan:  make(chan struct{}),
		parallel:  1,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Scheduler) Serve() error {
	s.taskChan = make(chan *Task, s.backlog)

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

type funcSchedulerOption func(*Scheduler)
type schedulerOption struct{}

var SchedulerOption schedulerOption

func (schedulerOption) WithBacklog(backlog int) funcSchedulerOption {
	return func(s *Scheduler) {
		s.WithBacklog(backlog)
	}
}

func (s *Scheduler) WithBacklog(backlog int) *Scheduler {
	s.backlog = backlog
	return s
}

func (schedulerOption) WithParallel(parallel int) funcSchedulerOption {
	return func(s *Scheduler) {
		s.WithParallel(parallel)
	}
}

func (s *Scheduler) WithParallel(parallel int) *Scheduler {
	s.parallel = parallel
	return s
}

func (schedulerOption) WithBackground() funcSchedulerOption {
	return func(s *Scheduler) {
		s.WithBackground()
	}
}

func (s *Scheduler) WithBackground() *Scheduler {
	s.background = true
	return s
}

func (schedulerOption) WithSafety() funcSchedulerOption {
	return func(s *Scheduler) {
		s.WithSafety()
	}
}

func (s *Scheduler) WithSafety() *Scheduler {
	s.safety = true
	return s
}

func (schedulerOption) WithErrorFunc(errorFunc func(error)) funcSchedulerOption {
	return func(s *Scheduler) {
		s.WithErrorFunc(errorFunc)
	}
}

func (s *Scheduler) WithErrorFunc(errorFunc func(error)) *Scheduler {
	s.errorFunc = errorFunc
	return s
}

func (schedulerOption) WithReportChan(reportChan chan<- *Report) funcSchedulerOption {
	return func(s *Scheduler) {
		s.WithReportChan(reportChan)
	}
}

func (s *Scheduler) WithReportChan(reportChan chan<- *Report) *Scheduler {
	s.reportChan = reportChan
	return s
}

func (schedulerOption) WithParallelChan(parallelChan <-chan int) funcSchedulerOption {
	return func(s *Scheduler) {
		s.WithParallelChan(parallelChan)
	}
}

func (s *Scheduler) WithParallelChan(parallelChan <-chan int) *Scheduler {
	s.parallelChan = parallelChan
	return s
}
