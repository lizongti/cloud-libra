package log

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/sirupsen/logrus"
)

var ErrHookNotFound = errors.New("hook not found")

type HookGenerator struct{}
type HookGenerateFunc func(*hierarchy.Hierarchy) (logrus.Hook, error)

var hookGeneratorMap = map[string]HookGenerateFunc{}

func init() {
	t := reflect.TypeOf(HookGenerator{})
	for index := 0; index < t.NumMethod(); index++ {
		method := t.Method(index)
		hookGeneratorMap[method.Name] = method.Func.Interface().(HookGenerateFunc)
	}
}

func NewHook(typ string, c *hierarchy.Hierarchy) (logrus.Hook, error) {
	if hookGenerator, ok := hookGeneratorMap[typ]; ok {
		return hookGenerator(c)
	}
	return nil, fmt.Errorf("%w: %s", ErrHookNotFound, typ)
}

func NewHooks(c *hierarchy.Hierarchy) ([]logrus.Hook, error) {
	hooks := make([]logrus.Hook, 0)

	c.ForeachInArray("hooks", func(index int, hierarchy *hierarchy.Hierarchy) (bool, error) {
		typ := hierarchy.GetString("type")
		hook, err := NewHook(typ, hierarchy)
		if err != nil {
			return false, err
		}
		hooks = append(hooks, hook)
		return true, nil
	})

	return hooks, nil
}
