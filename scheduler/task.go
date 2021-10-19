package scheduler

import (
	"log"
	"runtime"

	"github.com/gofrs/uuid"
)

type Task struct {
	taskOpt
	id        string
	doFunc    func()
	progress  int
	scheduler *Scheduler
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

func (t *Task) Progress() int {
	return t.progress
}

func (t *Task) MoveProgress(amount int) {
	t.progress += amount
	t.scheduler.report(t, RUNNING)
}

func (t *Task) Publish(s *Scheduler) {
	t.scheduler = s
	t.init()
	t.scheduler.schedule(t)
	t.scheduler.report(t, PENDING)
}

func (t *Task) execute() {
	defer func() {
		t.scheduler.report(t, FAILED)
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("scheduler: panic executing task %s: %v\n%s", t.name, err, buf)
		}
	}()
	t.scheduler.report(t, RUNNING)
	t.doFunc()
	t.scheduler.report(t, DONE)
}

func (t *Task) init() {
	uuid, _ := uuid.NewV4()
	t.id = uuid.String()
	if t.progressMax == 0 {
		t.progressMax = 1
	}
}

type taskOption func(*Task)
type taskOptions []taskOption

type taskOpt struct {
	taskOptions
	name        string
	progressMax int
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

func WithProgressMax(progressMax int) taskOption {
	return func(t *Task) {
		t.WithProgressMax(progressMax)
	}
}

func (t *Task) WithProgressMax(progressMax int) *Task {
	t.progressMax = progressMax
	return t
}
