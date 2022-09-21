package log

import (
	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/cloudlibraries/libra/log/hook"
	"github.com/sirupsen/logrus"
)

func NewHooks(c *hierarchy.Hierarchy) ([]logrus.Hook, error) {
	hooks := make([]logrus.Hook, 0)

	c.ForeachInArray("hooks", func(index int, hierarchy *hierarchy.Hierarchy) (bool, error) {
		typ := hierarchy.GetString("type")
		hook, err := hook.NewHook(typ, hierarchy)
		if err != nil {
			return false, err
		}
		hooks = append(hooks, hook)
		return true, nil
	})

	return hooks, nil
}
