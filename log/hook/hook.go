package hook

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/sirupsen/logrus"
)

var ErrHookNotFound = errors.New("hook not found")

type (
	Generator    struct{}
	GenerateFunc func(*hierarchy.Hierarchy) (logrus.Hook, error)
)

var hookGeneratorMap = map[string]GenerateFunc{}

func init() {
	t := reflect.TypeOf(&Generator{})
	for index := 0; index < t.NumMethod(); index++ {
		method := t.Method(index)
		hookGeneratorMap[strings.ToLower(method.Name)] = method.Func.Interface().(GenerateFunc)
	}
}

func NewHook(typ string, c *hierarchy.Hierarchy) (logrus.Hook, error) {
	if hookGenerator, ok := hookGeneratorMap[typ]; ok {
		return hookGenerator(c)
	}

	return nil, fmt.Errorf("%w: %s", ErrHookNotFound, typ)
}
