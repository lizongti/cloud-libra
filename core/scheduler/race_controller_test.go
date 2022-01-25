package scheduler_test

import (
	"testing"
	"time"

	"github.com/aceaura/libra/core/scheduler"
)

func TestRaceController(t *testing.T) {
	tasks := make([]*scheduler.Task, 0)
	tasks = append(tasks, scheduler.NewTask().WithStage(func(task *scheduler.Task) error {
		time.Sleep(1 * time.Second)
		return nil
	}), scheduler.NewTask().WithStage(func(task *scheduler.Task) error {
		time.Sleep(1 * time.Second)
		return nil
	}))

	c := scheduler.NewRaceController().WithSafety()
	c.WithDoneFunc(func(task *scheduler.Task) {
		t.Logf("%v done", task)
	})
	c.WithFailedFunc(func(task *scheduler.Task) {
		t.Logf("%v failed", task)
	})
	c.WithTask(tasks...)
	c.Serve()
}
