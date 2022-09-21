package formatter

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/cloudlibraries/libra/hierarchy"
	"github.com/sirupsen/logrus"
)

var (
	ErrFormatterNotFound = errors.New("formatter not found")
	ErrMethodNotValid    = errors.New("method not valid")
)

type (
	Generator    struct{}
	GenerateFunc func(*hierarchy.Hierarchy) (logrus.Formatter, error)
)

var formatterGeneratorMap = map[string]GenerateFunc{}

func init() {
	i := &Generator{}
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	for index := 0; index < t.NumMethod(); index++ {
		method := t.Method(index)
		formatterGeneratorMap[strings.ToLower(method.Name)] = func(h *hierarchy.Hierarchy) (logrus.Formatter, error) {
			in := []reflect.Value{v, reflect.ValueOf(h)}
			out := method.Func.Call(in)

			return out[0].Interface().(logrus.Formatter), out[1].Interface().(error)
		}
	}
}

func NewFormatter(c *hierarchy.Hierarchy) (logrus.Formatter, error) {
	typ := c.GetString("type")
	if fn, ok := formatterGeneratorMap[typ]; ok {
		return fn(c)
	}

	return nil, fmt.Errorf("%w: %s", ErrFormatterNotFound, typ)
}
