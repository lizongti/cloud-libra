package scheduler

import (
	"log"
	"runtime"

	"github.com/gofrs/uuid"
)

type Task struct {
	taskOpt
	scheduler *Scheduler
	id        string
	state     int
	progress  int
}

type Report struct {
	ID            string
	Name          string
	Progress      int
	TotalProgress int
	State         int
}

func NewTask(options ...taskOption) *Task {
	s := &Task{
		taskOpt: *newTaskOpt(options),
	}
	s.doOpt(s)
	return s
}

func (t *Task) String() string {
	return t.name
}

func (t *Task) Publish(s *Scheduler) {
	t.scheduler = s
	t.init()
	t.scheduler.schedule(t)

}

func (t *Task) init() {
	uuid, _ := uuid.NewV4()
	t.id = uuid.String()
}

func (t *Task) execute() {
	defer func() {
		if err := recover(); err != nil {
			t.switchState(FAILED)
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("scheduler: panic executing task %s: %v\n%s", t.name, err, buf)
		}
	}()
	t.switchState(RUNNING)

	for _, stage := range t.stages {
		if err := stage(); err != nil {
			panic(err)
		}
		t.switchStage()
	}

	t.switchState(DONE)
}

func (t *Task) switchState(state int) {
	t.state = state
	if t.scheduler.needReport() {
		t.scheduler.report(&Report{
			ID:            t.id,
			Name:          t.name,
			State:         t.state,
			Progress:      t.progress,
			TotalProgress: len(t.stages),
		})
	}
}

func (t *Task) switchStage() {
	t.progress++
	if t.scheduler.needReport() {
		t.scheduler.report(&Report{
			ID:            t.id,
			Name:          t.name,
			State:         t.state,
			Progress:      t.progress,
			TotalProgress: len(t.stages),
		})
	}
}

type taskOption func(*Task)
type taskOptions []taskOption

type taskOpt struct {
	taskOptions
	name   string
	stages []func() error
}

func newTaskOpt(options []taskOption) *taskOpt {
	return &taskOpt{
		taskOptions: options,
	}
}

func (opt *taskOpt) doOpt(s *Task) {
	for _, option := range opt.taskOptions {
		option(s)
	}
}

func WithName(name string) taskOption {
	return func(t *Task) {
		t.WithName(name)
	}
}

func (t *Task) WithName(name string) *Task {
	t.name = name
	return t
}

func WithStage(stage func() error) taskOption {
	return func(t *Task) {
		t.WithStage(stage)
	}
}

func (t *Task) WithStage(stage func() error) {
	t.stages = append(t.stages, stage)
}
