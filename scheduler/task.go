package scheduler

import (
	"fmt"
	"log"
	"runtime"

	"github.com/gofrs/uuid"
)

type Task struct {
	taskOpt
	scheduler     *Scheduler
	id            string
	state         int
	progress      int
	totalProgress int
	params        map[interface{}]interface{}
}

type Report struct {
	ID            string
	Name          string
	Progress      int
	TotalProgress int
	State         int
}

func NewTask(options ...taskOption) *Task {
	t := &Task{
		taskOpt: *newTaskOpt(options),
		params:  make(map[interface{}]interface{}),
	}
	t.doOpt(t)
	return t
}

func (t *Task) String() string {
	return fmt.Sprintf("%s[%s](%d/%d)", t.name, t.id, t.progress, t.totalProgress)
}

func (t *Task) ID() string {
	return t.id
}

func (t *Task) Name() string {
	return t.name
}

func (t *Task) Progress() int {
	return t.progress
}

func (t *Task) TotalProgress() int {
	return t.totalProgress
}

func (t *Task) State() int {
	return t.state
}

func (t *Task) Param(key interface{}) interface{} {
	value, _ := t.params[key]
	return value
}

func (t *Task) ParamBool(key interface{}) bool {
	return t.Param(key).(bool)
}

func (t *Task) ParamInt(key interface{}) int {
	return t.Param(key).(int)
}

func (t *Task) ParamString(key interface{}) string {
	return t.Param(key).(string)
}

func (t *Task) SetParam(key interface{}, value interface{}) {
	t.params[key] = value
}

func (t *Task) Publish(s *Scheduler) {
	t.scheduler = s
	t.init()
	t.switchState(TaskStatePending)
	t.scheduler.schedule(t)
}

func (t *Task) init() {
	uuid, _ := uuid.NewV4()
	t.id = uuid.String()
	t.state = TaskStateCreated
	t.totalProgress = len(t.stages)
}

func (t *Task) execute() {
	defer func() {
		if err := recover(); err != nil {
			t.switchState(TaskStateFailed)
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("scheduler: panic executing task %s: %v\n%s", t.name, err, buf)
		}
	}()
	t.switchState(TaskStateRunning)

	for _, stage := range t.stages {
		if err := stage(t); err != nil {
			panic(err)
		}
		t.switchStage()
	}

	t.switchState(TaskStateDone)
}

func (t *Task) switchState(state int) {
	t.state = state
	if t.scheduler.needReport() {
		t.scheduler.report(&Report{
			ID:            t.id,
			Name:          t.name,
			State:         t.state,
			Progress:      t.progress,
			TotalProgress: t.totalProgress,
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
			TotalProgress: t.totalProgress,
		})
	}
}

type taskOption func(*Task)
type taskOptions []taskOption

type taskOpt struct {
	taskOptions
	name   string
	stages []func(*Task) error
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

func WithStages(stages ...func(*Task) error) taskOption {
	return func(t *Task) {
		t.WithStages(stages...)
	}
}

func (t *Task) WithStages(stages ...func(*Task) error) *Task {
	t.stages = append(t.stages, stages...)
	return t
}

func WithParam(key interface{}, value interface{}) taskOption {
	return func(t *Task) {
		t.WithParam(key, value)
	}
}

func (t *Task) WithParam(key interface{}, value interface{}) *Task {
	t.params[key] = value
	return t
}
