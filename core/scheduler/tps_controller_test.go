package scheduler_test

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aceaura/libra/core/scheduler"
)

func TestTPSController(t *testing.T) {
	const (
		reportChanBacklog = 1000
		backlog           = 1000
		parallel          = 1
		parallelIncrease  = 1
		taskCount         = 100
		timeout           = 10
		parallelTick      = 100 * time.Millisecond
		tpsLimit          = 20
		taskBacklog       = 1000
		reportBacklog     = 4000
		parallelBacklog   = 1
	)
	var (
		errorChan = make(chan error, 1000)
	)
	c := scheduler.NewTPSController(
		scheduler.WithTPSSafety(),
		scheduler.WithTPSBackground(),
		scheduler.WithTPSErrorChan(errorChan),
		scheduler.WithTPSParallel(parallel),
		scheduler.WithTPSTaskBacklog(taskBacklog),
		scheduler.WithTPSReportBacklog(reportBacklog),
		scheduler.WithTPSParallelBacklog(parallelBacklog),
		scheduler.WithTPSParallelTick(parallelTick),
		scheduler.WithTPSParallelIncrease(parallelIncrease),
		scheduler.WithTPSLimit(tpsLimit),
	)

	if err := c.Serve(); err != nil {
		t.Fatal(err)
	}

	var count int64
	var exitChan = make(chan struct{}, 1)
	for index := 0; index < taskCount; index++ {
		scheduler.NewTask(
			scheduler.WithTaskName(fmt.Sprintf("test_parallel_tps[%d]", index)),
			scheduler.WithTaskStage(func(task *scheduler.Task) error {
				time.Sleep(time.Second * 1)
				atomic.AddInt64(&count, 1)
				if atomic.LoadInt64(&count) == taskCount {
					exitChan <- struct{}{}
				}
				return nil
			}),
		).Publish(c.Scheduler())
	}
	<-exitChan
}
