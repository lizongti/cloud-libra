package scheduler

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/aceaura/libra/magic"
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
	opts taskOptions

	id           string
	values       sync.Map
	state        TaskStateType
	progress     int
	taskStarted  time.Time
	stateStarted time.Time
	stageStarted time.Time
	context      context.Context
	report       func(r *Report)
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

func NewTask(opt ...ApplyTaskOption) *Task {
	uuid, _ := uuid.NewV4()
	now := time.Now()

	opts := defaultTaskOptions

	for _, o := range opt {
		o.apply(&opts)
	}

	t := &Task{
		opts:         opts,
		id:           uuid.String(),
		values:       sync.Map{},
		state:        TaskStateCreated,
		taskStarted:  now,
		stateStarted: now,
	}

	return t
}

func (t *Task) String() string {
	return fmt.Sprintf("%s[%s](%d/%d)", t.opts.name, t.id, t.progress, len(t.opts.stages))
}

func (t *Task) Context() context.Context {
	return t.context
}

func (t *Task) ID() string {
	return t.id
}

func (t *Task) Progress() int {
	return t.progress
}

func (t *Task) TotalProgress() int {
	return len(t.opts.stages)
}

func (t *Task) State() TaskStateType {
	return t.state
}

func (t *Task) Set(key interface{}, value interface{}) {
	t.values.Store(key, value)
}

func (t *Task) Value(key interface{}) interface{} {
	value, _ := t.values.Load(key)
	return value
}

func (t *Task) Publish(s *Scheduler) {
	t.report = func(r *Report) {
		s.report(r)
	}
	t.switchState(TaskStatePending)
	s.schedule(t)
}

func (t *Task) execute() {
	defer func() {
		if err := recover(); err != nil {
			t.switchState(TaskStateFailed)
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("scheduler: panic executing task [%s]: %v\n%s", t.opts.name, err, buf)
		}
	}()

	if t.opts.parentContext == nil {
		t.opts.parentContext = context.Background()
	}

	for k, v := range t.opts.params {
		t.Set(k, v)
	}

	if t.opts.timeout == 0 {
		ctx, cancel := context.WithCancel(t.opts.parentContext)
		defer cancel()

		t.context = ctx
		t.doStages()
		return
	}

	ctx, cancel := context.WithTimeout(t.opts.parentContext, t.opts.timeout)
	t.context = ctx
	doneChan := make(chan struct{})
	defer close(doneChan)
	defer cancel()

	go func() {
		t.doStages()
		doneChan <- struct{}{}
	}()

	select {
	case <-doneChan:
	case <-ctx.Done():
		panic(ctx.Err())
	}
}

func (t *Task) doStages() {
	t.switchState(TaskStateRunning)

	for _, stage := range t.opts.stages {
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

	if t.report != nil {
		t.report(&Report{
			ID:            t.id,
			Name:          t.opts.name,
			State:         t.state,
			Progress:      t.progress,
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

	if t.report != nil {
		t.report(&Report{
			ID:            t.id,
			Name:          t.opts.name,
			State:         t.state,
			Progress:      t.progress,
			TotalProgress: len(t.opts.stages),
			TaskDuration:  now.Sub(t.taskStarted),
			StageDuration: now.Sub(t.stageStarted),
			StateDuration: now.Sub(t.stateStarted),
		})
	}

	t.stageStarted = time.Now()
}

type taskOptions struct {
	name          string
	stages        []func(*Task) error
	params        map[interface{}]interface{}
	parentContext context.Context
	timeout       time.Duration
}

var defaultTaskOptions = taskOptions{
	name:          magic.Anonymous,
	stages:        nil,
	params:        nil,
	parentContext: nil,
	timeout:       0,
}

type ApplyTaskOption interface {
	apply(*taskOptions)
}

type funcTaskOption func(*taskOptions)

func (fto funcTaskOption) apply(to *taskOptions) {
	fto(to)
}

type taskOption int

var TaskOption taskOption

func (taskOption) Name(name string) funcTaskOption {
	return func(t *taskOptions) {
		t.name = name
	}
}

func (t *Task) WithName(name string) *Task {
	TaskOption.Name(name).apply(&t.opts)
	return t
}

func (taskOption) Stage(stages ...func(*Task) error) funcTaskOption {
	return func(t *taskOptions) {
		t.stages = append(t.stages, stages...)

	}
}

func (t *Task) WithStage(stages ...func(*Task) error) *Task {
	TaskOption.Stage(stages...).apply(&t.opts)
	return t
}

func (taskOption) Params(params map[interface{}]interface{}) funcTaskOption {
	return func(t *taskOptions) {
		t.params = params
	}
}

func (t *Task) WithParams(params map[interface{}]interface{}) *Task {
	TaskOption.Params(params).apply(&t.opts)
	return t
}

func (taskOption) ParentContext(context context.Context) funcTaskOption {
	return func(t *taskOptions) {
		t.parentContext = context
	}
}

func (t *Task) WithParentContext(context context.Context) *Task {
	TaskOption.ParentContext(context).apply(&t.opts)
	return t
}

func (taskOption) Timeout(timeout time.Duration) funcTaskOption {
	return func(t *taskOptions) {
		t.timeout = timeout
	}
}

func (t *Task) WithTimeout(timeout time.Duration) *Task {
	TaskOption.Timeout(timeout).apply(&t.opts)
	return t
}
