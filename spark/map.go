package spark

import (
	"github.com/aceaura/libra/boost/bar"
	"github.com/aceaura/libra/core/scheduler"
)

func Map(inputs []interface{}, f func(interface{}) (interface{}, error)) []interface{} {
	bar := bar.NewBar(len(inputs))
	bar.Begin()

	outputs := make([]interface{}, 0, len(inputs))
	round := 0
	for len(inputs) > 0 {
		round++
		var tasks []*scheduler.Task
		for _, input := range inputs {
			tasks = append(tasks, scheduler.NewTask(
				scheduler.WithTaskParams(map[interface{}]interface{}{"Input": input}),
				scheduler.WithTaskStage(func(task *scheduler.Task) error {
					output, err := f(task.Get("Input"))
					if err != nil {
						return err
					}
					task.Set("Output", output)
					return nil
				}),
			))
		}
		inputs = []interface{}{}
		c := scheduler.NewRaceController(
			scheduler.RaceControllerOption.WithSafety(),
			scheduler.RaceControllerOption.WithDoneFunc(func(task *scheduler.Task) {
				bar.Move(1)
				outputs = append(outputs, task.Get("Output"))
			}),
			scheduler.RaceControllerOption.WithFailedFunc(func(task *scheduler.Task) {
				inputs = append(inputs, task.Get("Input"))
			}),
			scheduler.RaceControllerOption.WithTasks(tasks...),
		)
		if err := c.Serve(); err != nil {
			panic(err)
		}
	}
	bar.End()

	return outputs
}

func MapTest(inputs []interface{}, f func(interface{}) (interface{}, error)) []interface{} {
	bar := bar.NewBar(len(inputs))
	bar.Begin()
	outputs := make([]interface{}, 0, len(inputs))
	for _, input := range inputs {
		output, err := f(input)
		if err != nil {
			panic(err)
		}
		outputs = append(outputs, output)
		bar.Move(1)
	}
	bar.End()
	return outputs
}
