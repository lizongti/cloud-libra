package tree

import (
	"strconv"

	"github.com/aceaura/libra/core/cast"
)

type MapTree struct {
	data map[string]interface{}
}

func NewMapTree(data map[string]interface{}) *MapTree {
	return &MapTree{
		data: data,
	}
}

func (mt *MapTree) Get(path []string) interface{} {
	return mt.search(mt.data, path)
}

func (mt *MapTree) search(source interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	switch source := source.(type) {
	case map[interface{}]interface{}:
		next, ok := cast.ToStringMap(source)[path[0]]
		if !ok {
			return nil
		}
		return mt.search(next, path[1:])
	case map[string]interface{}:
		next, ok := source[path[0]]
		if !ok {
			return nil
		}
		return mt.search(next, path[1:])
	case []interface{}:
		index, err := strconv.Atoi(path[0])
		if err != nil || len(source) <= index {
			return nil
		}
		next := source[index]
		return mt.search(next, path[1:])
	default:
		return nil
	}
}

func (mt *MapTree) Override(omt *MapTree) *MapTree {
	return nil
}
