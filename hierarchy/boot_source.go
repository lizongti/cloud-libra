package hierarchy

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

var ErrArgsNotEnough = errors.New("args not enough")

type SourceFunc func(string) (map[string][]byte, error)

type PlainSource struct{}

var plainSourceMap = map[string]SourceFunc{}

func init() {
	t := reflect.TypeOf(&PlainSource{})
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		plainSourceMap[strings.ToLower(method.Name)] = method.Func.Interface().(SourceFunc)
	}
}

func (*PlainSource) Plain(s string) (map[string][]byte, error) {
	assets := make(map[string][]byte)
	assets[""] = []byte(s)

	return assets, nil
}

type ModeSource struct {
	hierarchy *Hierarchy
}

var modeSourceMap = map[string]SourceFunc{}

func init() {
	t := reflect.TypeOf(&ModeSource{})
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		modeSourceMap[strings.ToLower(method.Name)] = method.Func.Interface().(SourceFunc)
	}
}

func (*ModeSource) Args(s string) (map[string][]byte, error) {
	assets := make(map[string][]byte)
	assets[""] = []byte(s)

	return assets, nil
}

func (*ModeSource) Env(s string) (map[string][]byte, error) {
	assets := make(map[string][]byte)
	assets[""] = []byte(s)

	return assets, nil
}

func (*ModeSource) Flags(s string) (map[string][]byte, error) {
	assets := make(map[string][]byte)
	assets[""] = []byte(s)

	return assets, nil
}

func (*ModeSource) Stdin(s string) (map[string][]byte, error) {
	data, err := io.ReadAll(bufio.NewReader(os.Stdin))
	if err != nil {
		return nil, err
	}

	assets := make(map[string][]byte)
	assets[""] = []byte(data)

	return assets, nil
}

func (ma *ModeSource) Hierarchy(s string) (map[string][]byte, error) {
	strs := strings.Split(s, ":")
	if len(strs) < 2 {
		return nil, fmt.Errorf("%w: %s", ErrArgsNotEnough, s)
	}
	key := strs[1]
	val := ma.hierarchy.GetString(key)

	assets := make(map[string][]byte)
	assets[""] = []byte(val)

	return assets, nil
}

type URLSource struct{}

var urlSourceMap = map[string]SourceFunc{}

func init() {
	t := reflect.TypeOf(&ModeSource{})
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		urlSourceMap[strings.ToLower(method.Name)] = method.Func.Interface().(SourceFunc)
	}
}

func (*URLSource) File(s string) (map[string][]byte, error) {
	panic("implement me")
}

func (*URLSource) Http(s string) (map[string][]byte, error) {
	panic("implement me")
}

func (*URLSource) Https(s string) (map[string][]byte, error) {
	panic("implement me")
}
