package scheduler

import (
	"fmt"
)

const (
	PENDING = iota
	RUNNING
	DONE
	FAILED
)

type Scheduler struct {
	schedulerOpt
	pipelines  []*pipeline
	taskChan   chan *Task
	reportChan chan *report
	panicChan  chan interface{}
	dieChan    chan struct{}
	exitChan   chan struct{}
	tablet     *Tablet
}

type pipeline struct {
	dieChan  chan struct{}
	exitChan chan struct{}
}

type report struct {
	id          string
	name        string
	progress    int
	progressMax int
	state       int
}

func NewScheduler() *Scheduler {
	s := &Scheduler{
		panicChan: make(chan interface{}, 1),
		dieChan:   make(chan struct{}),
		exitChan:  make(chan struct{}),
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
	case r := <-s.reportChan:
		if s.tablet != nil {
			s.tablet.readReport(r)
		}
	case p := <-s.panicChan:
		panic(p)
	case <-s.dieChan:
		return
	}

	return
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

func (s *Scheduler) init() {
	s.taskChan = make(chan *Task, s.taskBacklog)
	s.reportChan = make(chan *report, s.reportBacklog)

	for i := 0; i < s.parallel; i++ {
		s.pipelines = append(s.pipelines, &pipeline{
			dieChan:  make(chan struct{}),
			exitChan: make(chan struct{}),
		})
	}
}

func (s *Scheduler) close(p *pipeline) {
	close(p.dieChan)
	<-p.exitChan
}

func (s *Scheduler) report(t *Task, status int) {
	if s.tablet != nil {
		s.reportChan <- s.tablet.newReport(t, status)
	}
}

type Tablet struct {
	PendingTasks map[int]*report
	DoneTasks    []*report
}

func (t *Tablet) readReport(r *report) {
	// TODO
}

func (*Tablet) newReport(t *Task, status int) *report {
	return &report{
		id:          t.id,
		name:        t.name,
		progress:    t.progress,
		progressMax: t.progressMax,
		state:       status,
	}
}

func (*Tablet) TasksByState(states ...int) []*report {
	// TODO
	return nil
}

type schedulerOption func(*Scheduler)
type schedulerOptions []schedulerOption

type schedulerOpt struct {
	schedulerOptions
	taskBacklog   int
	reportBacklog int
	tpsLimit      int
	parallel      int
	background    bool
	safety        bool
	errorFunc     func(error)
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

func WithTaskBacklog(taskBacklog int) schedulerOption {
	return func(s *Scheduler) {
		s.WithTaskBacklog(taskBacklog)
	}
}

func (s *Scheduler) WithTaskBacklog(taskBacklog int) *Scheduler {
	s.taskBacklog = taskBacklog
	return s
}

func WithReportBacklog(reportBacklog int) schedulerOption {
	return func(s *Scheduler) {
		s.WithReportBacklog(reportBacklog)
	}
}

func (s *Scheduler) WithReportBacklog(reportBacklog int) *Scheduler {
	s.reportBacklog = reportBacklog
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
