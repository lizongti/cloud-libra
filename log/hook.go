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
	c.GetStringMapStringSlice()
	hooks := make([]logrus.Hook, 0)
	for _, typ := range  {
		hook, err := NewHook(typ, c)
		if err != nil {
			return nil, err
		}
		hooks = append(hooks, hook)
	}
	return hooks, nil
}