package scheduler

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/gofrs/uuid"
)

type TaskStateType int

const (
	TaskStateCreated TaskStateType = iota
	TaskStatePending
	TaskStateRunning
	TaskStateDone
	TaskStateFailed
)

var taskStateName = map[TaskStateType]string{
	TaskStateCreated: "created",
	TaskStatePending: "pending",
	TaskStateRunning: "running",
	TaskStateDone:    "done",
	TaskStateFailed:  "failed",
}

func (t TaskStateType) String() string {
	if s, ok := taskStateName[t]; ok {
		return s
	}
	return fmt.Sprintf("taskStateName=%d?", int(t))
}

type Task struct {
	id            string
	name          string
	stages        []func(*Task) error
	params        map[interface{}]interface{}
	context       context.Context
	timeout       time.Duration
	scheduler     *Scheduler
	state         TaskStateType
	progress      int
	totalProgress int
	taskStarted   time.Time
	stateStarted  time.Time
	stageStarted  time.Time
}

type Report struct {
	ID            string
	Name          string
	Progress      int
	TotalProgress int
	State         TaskStateType
	TaskDuration  time.Duration
	StateDuration time.Duration
	StageDuration time.Duration
}

func NewTask(opts ...taskOpt) *Task {
	uuid, _ := uuid.NewV4()
	now := time.Now()
	t := &Task{
		id:           uuid.String(),
		name:         "anonymous",
		state:        TaskStateCreated,
		params:       make(map[interface{}]interface{}),
		taskStarted:  now,
		stateStarted: now,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

func (t *Task) String() string {
	return fmt.Sprintf("%s[%s](%d/%d)", t.name, t.id, t.progress, t.totalProgress)
}

func (t *Task) Context() context.Context {
	return t.context
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

func (t *Task) State() TaskStateType {
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
	t.switchState(TaskStatePending)
	t.scheduler.schedule(t)
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

func (t *Task) switchState(state TaskStateType) {
	t.state = state
	now := time.Now()
	if t.scheduler.needReport() {
		t.scheduler.report(&Report{
			ID:            t.id,
			Name:          t.name,
			State:         t.state,
			Progress:      t.progress,
			TotalProgress: t.totalProgress,
			TaskDuration:  now.Sub(t.taskStarted),
			StageDuration: 0,
			StateDuration: now.Sub(t.stateStarted),
		})
	}
	if t.state == TaskStateRunning {
		t.stageStarted = now
	} else {
		t.stageStarted = time.Time{}
	}

	t.stateStarted = now
}

func (t *Task) switchStage() {
	t.progress++
	now := time.Now()
	if t.scheduler.needReport() {
		t.scheduler.report(&Report{
			ID:            t.id,
			Name:          t.name,
			State:         t.state,
			Progress:      t.progress,
			TotalProgress: t.totalProgress,
			TaskDuration:  now.Sub(t.taskStarted),
			StageDuration: now.Sub(t.stageStarted),
			StateDuration: now.Sub(t.stateStarted),
		})
	}
	t.stageStarted = time.Now()
}

type taskOpt func(*Task)
type taskOption struct{}

var TaskOption taskOption

func (taskOption) WithName(name string) taskOpt {
	return func(t *Task) {
		t.WithName(name)
	}
}

func (t *Task) WithName(name string) *Task {
	t.name = name
	return t
}

func (taskOption) WithStage(stages ...func(*Task) error) taskOpt {
	return func(t *Task) {
		t.WithStage(stages...)
	}
}

func (t *Task) WithStage(stages ...func(*Task) error) *Task {
	t.stages = append(t.stages, stages...)
	t.totalProgress = len(t.stages)
	return t
}

func (taskOption) WithParam(key interface{}, value interface{}) taskOpt {
	return func(t *Task) {
		t.WithParam(key, value)
	}
}

func (t *Task) WithParam(key interface{}, value interface{}) *Task {
	t.params[key] = value
	return t
}

func (taskOption) WithParams(params map[interface{}]interface{}) taskOpt {
	return func(t *Task) {
		t.WithParams(params)
	}
}

func (t *Task) WithParams(params map[interface{}]interface{}) *Task {
	t.params = params
	return t
}

func (taskOption) WithContext(context context.Context) taskOpt {
	return func(t *Task) {
		t.WithTask(context)
	}
}

func (t *Task) WithTask(context context.Context) *Task {
	t.context = context
	return t
}

func (taskOption) WithTimeout(timeout time.Duration) taskOpt {
	return func(t *Task) {
		t.WithTimeout(timeout)
	}
}

func (t *Task) WithTimeout(timeout time.Duration) *Task {
	t.timeout = timeout
	return t
}
