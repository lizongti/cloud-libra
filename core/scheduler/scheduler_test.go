package scheduler_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aceaura/libra/core/scheduler"
)

func TestParallel(t *testing.T) {
	const (
		reportChanBacklog = 1000
		backlog           = 1000
		parallel          = 1000
		taskCount         = 300
		timeout           = 2
	)
	var reportChan = make(chan *scheduler.Report, reportChanBacklog)
	s := scheduler.NewScheduler().WithReportChan(reportChan).WithParallel(parallel).WithTaskBacklog(backlog)
	if err := s.WithBackground().Serve(); err != nil {
		t.Fatalf("unexpected error getting from scheduler: %v", err)
	}
	for index := 0; index < taskCount; index++ {
		scheduler.NewTask(
			scheduler.TaskOption.Name(fmt.Sprintf("test_parallel[%d]", index)),
			scheduler.TaskOption.Stage(func(task *scheduler.Task) error {
				time.Sleep(time.Second * 1)
				return nil
			}),
		).Publish(s)
	}
	var finishCount = 0
	var timeoutChan = time.After(time.Duration(timeout) * time.Second)
	for {
		select {
		case <-timeoutChan:
			t.Fatal("timed out when getting report from task")
		case r := <-reportChan:
			if r.State == scheduler.TaskStateDone || r.State == scheduler.TaskStateFailed {
				finishCount++
				if finishCount == taskCount {
					return
				}
			}
		}
	}
}

func TestParallelChan(t *testing.T) {
	const (
		reportChanBacklog = 1000
		backlog           = 1000
		parallel          = 1
		parallelIncrease  = 1
		taskCount         = 300
		timeout           = 10
		parallelTickMS    = 1
	)
	var reportChan = make(chan *scheduler.Report, reportChanBacklog)
	var parallelChan = make(chan int)
	s := scheduler.NewScheduler(
		scheduler.SchedulerOption.TaskBacklog(backlog),
		scheduler.SchedulerOption.Parallel(parallel),
		scheduler.SchedulerOption.ReportChan(reportChan),
		scheduler.SchedulerOption.ParallelChan(parallelChan),
	)
	if err := s.WithBackground().Serve(); err != nil {
		t.Fatalf("unexpected error getting from scheduler: %v", err)
	}
	for index := 0; index < taskCount; index++ {
		scheduler.NewTask(
			scheduler.TaskOption.Name(fmt.Sprintf("test_parallel_chan[%d]", index)),
			scheduler.TaskOption.Stage(func(task *scheduler.Task) error {
				time.Sleep(time.Second * 1)
				return nil
			}),
		).Publish(s)
		t.Log(index)
	}
	var timeoutChan = time.After(time.Duration(timeout) * time.Second)
	var tickerChan = time.NewTicker(time.Duration(parallelTickMS) * time.Millisecond).C
	var stateMap = make(map[scheduler.TaskStateType]int)
	var finishCount int
	for {
		select {
		case <-timeoutChan:
			t.Fatal("timeout when getting report from task")
		case r := <-reportChan:
			if r.Progress == 0 {
				stateMap[r.State]++
			}
			if r.State == scheduler.TaskStateDone || r.State == scheduler.TaskStateFailed {
				finishCount++
				if finishCount == taskCount {
					return
				}
			}
		case <-tickerChan:
			if stateMap[scheduler.TaskStatePending]-stateMap[scheduler.TaskStateRunning] > 0 {
				parallelChan <- parallelIncrease
			}
		}
	}
}

func TestTPS(t *testing.T) {
	const (
		reportChanBacklog = 1000
		backlog           = 1000
		parallel          = 1
		parallelIncrease  = 1
		taskCount         = 100
		timeout           = 10
		parallelTickMS    = 100
		tps               = 20
	)
	var reportChan = make(chan *scheduler.Report, reportChanBacklog)
	var parallelChan = make(chan int)
	s := scheduler.NewScheduler(
		scheduler.SchedulerOption.TaskBacklog(backlog),
		scheduler.SchedulerOption.Parallel(parallel),
		scheduler.SchedulerOption.ReportChan(reportChan),
		scheduler.SchedulerOption.ParallelChan(parallelChan),
	)
	if err := s.WithBackground().Serve(); err != nil {
		t.Fatalf("unexpected error getting from scheduler: %v", err)
	}
	for index := 0; index < taskCount; index++ {
		scheduler.NewTask(
			scheduler.TaskOption.Name(fmt.Sprintf("test_parallel_tps[%d]", index)),
			scheduler.TaskOption.Stage(func(task *scheduler.Task) error {
				time.Sleep(time.Second * 1)
				return nil
			}),
		).Publish(s)
	}
	var timeoutChan = time.After(time.Duration(timeout) * time.Second)
	var tickerChan = time.NewTicker(time.Duration(parallelTickMS) * time.Millisecond).C
	var tpsTickerChan = time.NewTicker(time.Second).C
	var stateMap = make(map[scheduler.TaskStateType]int)
	var finishCount int
	var tpsFinishCount int
	var maxTPS int
	for {
		select {
		case <-timeoutChan:
			t.Fatal("timed out when getting report from task")
		case r := <-reportChan:
			if r.Progress == 0 {
				stateMap[r.State]++
			}
			if r.State == scheduler.TaskStateDone || r.State == scheduler.TaskStateFailed {
				finishCount++
				tpsFinishCount++
				if finishCount == taskCount {
					return
				}
			}
		case <-tickerChan:
			if stateMap[scheduler.TaskStatePending]-stateMap[scheduler.TaskStateRunning] > 0 {
				if maxTPS < tps {
					parallelChan <- parallelIncrease
				}
			}
		case <-tpsTickerChan:
			if tpsFinishCount > maxTPS {
				maxTPS = tpsFinishCount
			}
			tpsFinishCount = 0
		}
	}
}
