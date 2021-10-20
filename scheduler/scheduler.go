package scheduler

import (
	"fmt"
)

const (
	TaskStateCreated = iota
	TaskStatePending
	TaskStateRunning
	TaskStateDone
	TaskStateFailed
)

type Scheduler struct {
	*schedulerOpt
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

func NewScheduler(options ...schedulerOption) *Scheduler {
	s := &Scheduler{
		schedulerOpt: newSchedulerOpt(options),
		panicChan:    make(chan interface{}, 1),
		dieChan:      make(chan struct{}),
		exitChan:     make(chan struct{}),
	}
	s.doOpt(s)
	return s
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
	s.taskChan = make(chan *Task, s.backlog)
	if s.parallel == 0 {
		s.parallel = 1
	}

	for i := 0; i < s.parallel; i++ {
		s.pipelines = append(s.pipelines, &pipeline{
			dieChan:  make(chan struct{}),
			exitChan: make(chan struct{}),
		})
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

	for _, p := range s.pipelines {
		go func(p *pipeline) {
			s.digest(p)
		}(p)
	}

	select {
	case p := <-s.panicChan:
		panic(p)
	case <-s.dieChan:
		return
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

type schedulerOption func(*Scheduler)
type schedulerOptions []schedulerOption

type schedulerOpt struct {
	schedulerOptions
	backlog    int
	tpsLimit   int
	parallel   int
	background bool
	safety     bool
	errorFunc  func(error)
	reportChan chan<- *Report
}

func newSchedulerOpt(options []schedulerOption) *schedulerOpt {
	return &schedulerOpt{
		schedulerOptions: options,
	}
}

func (opt *schedulerOpt) doOpt(s *Scheduler) {
	for _, option := range opt.schedulerOptions {
		option(s)
	}
}

func WithBacklog(taskBacklog int) schedulerOption {
	return func(s *Scheduler) {
		s.WithBacklog(taskBacklog)
	}
}

func (s *Scheduler) WithBacklog(taskBacklog int) *Scheduler {
	s.backlog = taskBacklog
	return s
}

func WithTPSLimit(tpsLimit int) schedulerOption {
	return func(s *Scheduler) {
		s.WithTPSLimit(tpsLimit)
	}
}

func (s *Scheduler) WithTPSLimit(tpsLimit int) *Scheduler {
	s.tpsLimit = tpsLimit
	return s
}

func WithParallel(parallel int) schedulerOption {
	return func(s *Scheduler) {
		s.WithParallel(parallel)
	}
}

func (s *Scheduler) WithParallel(parallel int) *Scheduler {
	s.parallel = parallel
	return s
}

func WithBackground() schedulerOption {
	return func(s *Scheduler) {
		s.WithBackground()
	}
}

func (s *Scheduler) WithBackground() *Scheduler {
	s.background = true
	return s
}

func WithSchdedulerSafety() schedulerOption {
	return func(s *Scheduler) {
		s.WithSchdedulerSafety()
	}
}

func (s *Scheduler) WithSchdedulerSafety() *Scheduler {
	s.safety = true
	return s
}

func WithErrorFunc(errorFunc func(error)) schedulerOption {
	return func(s *Scheduler) {
		s.WithErrorFunc(errorFunc)
	}
}

func (s *Scheduler) WithErrorFunc(errorFunc func(error)) *Scheduler {
	s.errorFunc = errorFunc
	return s
}

func WithReportChan(reportChan chan<- *Report) schedulerOption {
	return func(s *Scheduler) {
		s.WithReportChan(reportChan)
	}
}

func (s *Scheduler) WithReportChan(reportChan chan<- *Report) *Scheduler {
	s.reportChan = reportChan
	return s
}
