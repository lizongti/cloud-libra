package scheduler_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aceaura/libra/scheduler"
)

func TestParallelPipeline(t *testing.T) {
	const (
		reportChanBacklog = 1000
		backlog           = 1000
		parallel          = 1000
		taskCount         = 300
		timeout           = 2
	)
	var reportChan = make(chan *scheduler.Report, reportChanBacklog)
	s := scheduler.NewScheduler().WithReportChan(reportChan).WithParallel(parallel).WithBacklog(backlog)
	if err := s.WithBackground().Serve(); err != nil {
		t.Fatalf("unexpected error getting from scheduler: %v", err)
	}
	for index := 0; index < taskCount; index++ {
		scheduler.NewTask(
			scheduler.WithName(fmt.Sprintf("test_parallel_pipeline[%d]", index)),
			scheduler.WithStages(func(task *scheduler.Task) error {
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
			t.Fatal("timeout when getting report from task")
		case r := <-reportChan:
			if r.State == scheduler.TaskStateDone {
				finishCount++
			}
			if finishCount == taskCount {
				return
			}
		}
	}
}
