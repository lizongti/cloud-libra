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
	if s != emptyScheduler && s.reportChan != nil {
		s.reportChan <- r
	}
}

type funcSchedulerOption func(*Scheduler)
type schedulerOption struct{}

var SchedulerOption schedulerOption

func (schedulerOption) Backlog(backlog int) funcSchedulerOption {
	return func(s *Scheduler) {
		s.Backlog(backlog)
	}
}

func (s *Scheduler) Backlog(backlog int) *Scheduler {
	s.backlog = backlog
	return s
}

func (schedulerOption) Parallel(parallel int) funcSchedulerOption {
	return func(s *Scheduler) {
		s.Parallel(parallel)
	}
}

func (s *Scheduler) Parallel(parallel int) *Scheduler {
	s.parallel = parallel
	return s
}

func (schedulerOption) Background() funcSchedulerOption {
	return func(s *Scheduler) {
		s.Background()
	}
}

func (s *Scheduler) Background() *Scheduler {
	s.background = true
	return s
}

func (schedulerOption) Safety() funcSchedulerOption {
	return func(s *Scheduler) {
		s.Safety()
	}
}

func (s *Scheduler) Safety() *Scheduler {
	s.safety = true
	return s
}

func (schedulerOption) WithErrorFunc(errorFunc func(error)) funcSchedulerOption {
	return func(s *Scheduler) {
		s.ErrorFunc(errorFunc)
	}
}

func (s *Scheduler) ErrorFunc(errorFunc func(error)) *Scheduler {
	s.errorFunc = errorFunc
	return s
}

func (schedulerOption) ReportChan(reportChan chan<- *Report) funcSchedulerOption {
	return func(s *Scheduler) {
		s.ReportChan(reportChan)
	}
}

func (s *Scheduler) ReportChan(reportChan chan<- *Report) *Scheduler {
	s.reportChan = reportChan
	return s
}

func (schedulerOption) ParallelChan(parallelChan <-chan int) funcSchedulerOption {
	return func(s *Scheduler) {
		s.ParallelChan(parallelChan)
	}
}

func (s *Scheduler) ParallelChan(parallelChan <-chan int) *Scheduler {
	s.parallelChan = parallelChan
	return s
}
