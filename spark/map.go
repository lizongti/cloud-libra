package spark

import (
	"github.com/aceaura/libra/boost/bar"
	"github.com/aceaura/libra/core/scheduler"
)

func Map(inputs []interface{}, f func(interface{}) interface{}) []interface{} {
	bar := bar.NewBar(len(inputs))

	outputs := make([]interface{}, 0, len(inputs))
	round := 0
	for len(inputs) > 0 {
		round++
		var tasks []*scheduler.Task
		for _, input := range inputs {
			tasks = append(tasks, scheduler.NewTask(
				scheduler.TaskOption.Params(map[interface{}]interface{}{"Input": input}),
				scheduler.TaskOption.Stage(func(task *scheduler.Task) error {
					task.Set("Output", f(task.Get("Input")))
					return nil
				}),
			))
		}
		c := scheduler.NewRaceController().WithSafety()
		c.WithDoneFunc(func(task *scheduler.Task) {
			bar.Move(1)
			outputs = append(outputs, task.Get("Output"))
		})
		inputs = []interface{}{}
		c.WithFailedFunc(func(task *scheduler.Task) {
			inputs = append(inputs, task.Get("Input"))
		})
		c.WithTask(tasks...)
		if err := c.Serve(); err != nil {
			panic(err)
		}
	}
	bar.Close()

	return outputs
}

func TestMap(inputs []interface{}, f func(interface{}) interface{}) []interface{} {
	bar := bar.NewBar(len(inputs))

	outputs := make([]interface{}, 0, len(inputs))
	for _, input := range inputs {
		outputs = append(outputs, f(input))
		bar.Move(1)
	}
	bar.Close()
	return outputs
}
