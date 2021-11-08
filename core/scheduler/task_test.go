package scheduler_test

import (
	"testing"
	"time"

	"github.com/aceaura/libra/scheduler"
)

func TestTaskState(t *testing.T) {
	const (
		reportChanBacklog = 1000
		timeout           = 1
	)
	var reportChan = make(chan *scheduler.Report, reportChanBacklog)
	s := scheduler.NewScheduler().WithReportChan(reportChan)
	if err := s.WithBackground().Serve(); err != nil {
		t.Fatalf("unexpected error getting from scheduler: %v", err)
	}
	scheduler.NewTask().WithName("test_task_state").Publish(s)
	var states = []scheduler.TaskStateType{
		scheduler.TaskStateCreated,
		scheduler.TaskStatePending,
		scheduler.TaskStateRunning,
		scheduler.TaskStateDone,
	}
	var stateIndex = 0
	var timeoutChan = time.After(time.Duration(timeout) * time.Second)
	for {
		select {
		case <-timeoutChan:
			t.Fatal("timeout when getting report from task")
		case r := <-reportChan:
			stateIndex++
			state := states[stateIndex]
			if state != r.State {
				t.Fatalf("expected a state of %d, got %d", state, r.State)
			}
			if state == scheduler.TaskStateDone {
				return
			}
		}
	}
}

func TestTaskStage(t *testing.T) {
	const (
		reportChanBacklog = 1000
		stageCount        = 100
		timeout           = 1
	)
	var reportChan = make(chan *scheduler.Report, reportChanBacklog)
	s := scheduler.NewScheduler().WithReportChan(reportChan)
	if err := s.WithBackground().Serve(); err != nil {
		t.Fatalf("unexpected error getting from scheduler: %v", err)
	}
	var stages = make([]func(*scheduler.Task) error, 0, stageCount)
	for index := 0; index < 10; index++ {
		stages = append(stages, func(task *scheduler.Task) error {
			return nil
		})
	}

	scheduler.NewTask().WithStage(stages...).WithName("test_task_stage").Publish(s)
	var progress = -1
	var timeoutChan = time.After(time.Duration(timeout) * time.Second)
	for {
		select {
		case <-timeoutChan:
			t.Fatal("timeout when getting report from task")
		case r := <-reportChan:
			if r.State == scheduler.TaskStateRunning {
				progress++
				if r.Progress != progress {
					t.Fatalf("getting error progress from report: %d", r.Progress)
				}
			}
			if r.State == scheduler.TaskStateDone {
				return
			}
		}
	}
}

func TestTaskParams(t *testing.T) {
	const (
		reportChanBacklog = 1000
		stageCount        = 100
		timeout           = 1
	)
	var reportChan = make(chan *scheduler.Report, reportChanBacklog)
	s := scheduler.NewScheduler().WithReportChan(reportChan)
	if err := s.WithBackground().Serve(); err != nil {
		t.Fatalf("unexpected error getting from scheduler: %v", err)
	}
	var stages = make([]func(*scheduler.Task) error, 0, stageCount)
	for index := 0; index < 10; index++ {
		stages = append(stages, func(task *scheduler.Task) error {
			progress := task.ParamInt("progress")
			defer func() {
				task.SetParam("progress", progress)
			}()

			if progress != task.Progress() {
				t.Fatalf("param progress not correct, expect %d, got %d",
					task.Progress(), progress)
			}

			progress++
			return nil
		})
	}

	scheduler.NewTask().WithStage(stages...).WithParam(
		"progress", 0).WithName("test_task_params").Publish(s)

	var timeoutChan = time.After(time.Duration(timeout) * time.Second)
	for {
		select {
		case <-timeoutChan:
			t.Fatal("timeout when getting report from task")
		case r := <-reportChan:
			if r.State == scheduler.TaskStateDone {
				return
			}
		}
	}
}

func TestTaskTimeout(t *testing.T) {
	const (
		reportChanBacklog = 1000
		timeout           = 3
		taskTimeout       = 1
		sleep             = 2
	)
	var reportChan = make(chan *scheduler.Report, reportChanBacklog)
	s := scheduler.NewScheduler().WithReportChan(reportChan)
	if err := s.WithBackground().Serve(); err != nil {
		t.Fatalf("unexpected error getting from scheduler: %v", err)
	}
	scheduler.NewTask().WithStage(func(*scheduler.Task) error {
		time.Sleep(time.Duration(sleep) * time.Second)
		return nil
	}).WithName("test_task_timeout").WithTimeout(time.Duration(taskTimeout) * time.Second).Publish(s)
	var timeoutChan = time.After(time.Duration(timeout) * time.Second)
	for {
		select {
		case <-timeoutChan:
			t.Fatal("timeout when getting report from task")
		case r := <-reportChan:
			if r.State == scheduler.TaskStateFailed {
				return
			}
		}
	}
}

func TestTaskReportTime(t *testing.T) {
	const (
		reportChanBacklog = 1000
		stageCount        = 100
		timeout           = 20
		sleep             = 1
	)
	var reportChan = make(chan *scheduler.Report, reportChanBacklog)
	s := scheduler.NewScheduler().WithReportChan(reportChan)
	if err := s.WithBackground().Serve(); err != nil {
		t.Fatalf("unexpected error getting from scheduler: %v", err)
	}
	var stages = make([]func(*scheduler.Task) error, 0, stageCount)
	for index := 0; index < 10; index++ {
		stages = append(stages, func(task *scheduler.Task) error {
			time.Sleep(time.Duration(sleep) * time.Second)
			return nil
		})
	}

	scheduler.NewTask().WithStage(stages...).WithName("test_task_report_time").Publish(s)
	var timeoutChan = time.After(time.Duration(timeout) * time.Second)
	for {
		select {
		case <-timeoutChan:
			t.Fatal("timeout when getting report from task")
		case r := <-reportChan:
			if r.State == scheduler.TaskStateRunning && r.Progress > 0 {
				if r.TaskDuration.Nanoseconds() == 0 {
					t.Fatal("expecting a task duration, got 0")
				}
				if r.StateDuration.Nanoseconds() == 0 {
					t.Fatal("expecting a state duration, got 0")
				}
				if r.StageDuration.Nanoseconds() == 0 {
					t.Fatal("expecting a stage duration, got 0")
				}
			}
			t.Logf("state: %v, progress: %v, stage task duration: %v(ms), state duration: %v(ms), stage duration: %v(ms)",
				r.State,
				r.Progress,
				r.TaskDuration.Milliseconds(),
				r.StateDuration.Milliseconds(),
				r.StageDuration.Milliseconds(),
			)
			if r.State == scheduler.TaskStateDone {
				return
			}
		}
	}
}
