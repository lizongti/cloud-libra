package scheduler_test

import (
	"testing"
	"time"

	"github.com/aceaura/libra/core/scheduler"
)

func TestRaceController(t *testing.T) {
	tasks := make([]*scheduler.Task, 0)
	for i := 0; i < 4000; i++ {
		tasks = append(tasks, scheduler.NewTask(
			scheduler.WithTaskStage(func(task *scheduler.Task) error {
				time.Sleep(1 * time.Second)
				return nil
			})))
	}

	c := scheduler.NewRaceController(
		scheduler.RaceControllerOption.WithSafety(),
		scheduler.RaceControllerOption.WithDoneFunc(func(task *scheduler.Task) {
			t.Logf("%v done", task)
		}),
		scheduler.RaceControllerOption.WithFailedFunc(func(task *scheduler.Task) {
			t.Logf("%v failed", task)
		}),
		scheduler.RaceControllerOption.WithTasks(tasks...),
	)
	c.Serve()
}
