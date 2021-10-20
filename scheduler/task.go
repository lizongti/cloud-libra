package scheduler

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/gofrs/uuid"
)

type Task struct {
	taskOpt
	scheduler      *Scheduler
	id             string
	state          int
	progress       int
	totalProgress  int
	taskBeginTime  int64
	stateBeginTime int64
	stageBeginTime int64
}

type Report struct {
	ID            string
	Name          string
	Progress      int
	TotalProgress int
	State         int
	TaskCost      int64
	LastStateCost int64
	LastStageCost int64
}

func NewTask(options ...taskOption) *Task {
	t := &Task{
		taskOpt: *newTaskOpt(options),
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

func (t *Task) SetParam(key interface{}, value interface{}) {
	t.params[key] = value
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
	if t.name == "" {
		t.name = "anonymous"
	}
	now := time.Now().UnixNano()
	t.taskBeginTime = now
	t.stateBeginTime = now
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

	if t.timeout == 0 {
		t.doStages()
		return
	}

	doneChan := make(chan struct{})

	go func() {
		t.doStages()
		doneChan <- struct{}{}
	}()

	select {
	case <-doneChan:
	case <-time.After(t.timeout):
		panic(fmt.Errorf("do stages timeout"))
	}
}

func (t *Task) doStages() {
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
	now := time.Now().UnixNano()
	if t.scheduler.needReport() {
		t.scheduler.report(&Report{
			ID:            t.id,
			Name:          t.name,
			State:         t.state,
			Progress:      t.progress,
			TotalProgress: t.totalProgress,
			LastStageCost: 0,
			LastStateCost: now - t.stateBeginTime,
		})
	}
	if t.state == TaskStateRunning {
		t.stageBeginTime = now
	} else {
		t.stageBeginTime = 0
	}

	t.stateBeginTime = now
}

func (t *Task) switchStage() {
	t.progress++
	now := time.Now().UnixNano()
	if t.scheduler.needReport() {
		t.scheduler.report(&Report{
			ID:            t.id,
			Name:          t.name,
			State:         t.state,
			Progress:      t.progress,
			TotalProgress: t.totalProgress,
			LastStageCost: now - t.stageBeginTime,
			LastStateCost: now - t.stateBeginTime,
		})
	}
	t.stageBeginTime = now
}

type taskOption func(*Task)
type taskOptions []taskOption

type taskOpt struct {
	taskOptions
	name    string
	stages  []func(*Task) error
	params  map[interface{}]interface{}
	timeout time.Duration
}

func newTaskOpt(options []taskOption) *taskOpt {
	return &taskOpt{
		taskOptions: options,
		params:      make(map[interface{}]interface{}),
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

func WithParams(params map[interface{}]interface{}) taskOption {
	return func(t *Task) {
		t.WithParams(params)
	}
}

func (t *Task) WithParams(params map[interface{}]interface{}) *Task {
	t.params = params
	return t
}

func WithTimeout(timeout time.Duration) taskOption {
	return func(t *Task) {
		t.WithTimeout(timeout)
	}
}

func (t *Task) WithTimeout(timeout time.Duration) *Task {
	t.timeout = timeout
	return t
}
