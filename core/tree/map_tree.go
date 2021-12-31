package tree

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/aceaura/libra/core/deepcopy"
	"github.com/aceaura/libra/core/magic"
)

type MapTree struct {
	data map[string]interface{}
}

func NewMapTree() *MapTree {
	return &MapTree{
		data: make(map[string]interface{}),
	}
}

func (mt *MapTree) SetData(data map[string]interface{}) {
	mt.data = data
}

func (mt *MapTree) Data() interface{} {
	return mt.data
}

func (mt *MapTree) Get(path []string) interface{} {
	return mt.get(mt.data, path)
}

func (mt *MapTree) get(source interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	switch source := source.(type) {
	case map[string]interface{}:
		next, ok := source[path[0]]
		if !ok {
			return nil
		}
		return mt.get(next, path[1:])
	case []interface{}:
		index, err := strconv.Atoi(path[0])
		if err != nil || len(source) <= index {
			return nil
		}
		next := source[index]
		return mt.get(next, path[1:])
	default:
		return nil
	}
}

func (mt *MapTree) Set(path []string, v interface{}) {
	mt.data = mt.set(mt.data, path, v).(map[string]interface{})
}

func (mt *MapTree) set(source interface{}, path []string, v interface{}) interface{} {
	if len(path) == 0 {
		return source
	}

	if len(path) == 1 {
		switch source := source.(type) {
		case map[string]interface{}:
			source[path[0]] = v
			return source
		case []interface{}:
			index, err := strconv.Atoi(path[0])
			if err != nil {
				return nil
			}
			for len(source) <= index {
				source = append(source, nil)
			}
			source[index] = v
			return source
		default:
			return source
		}
	}

	switch source := source.(type) {
	case map[string]interface{}:
		next, ok := source[path[0]]
		if !ok {
			_, err := strconv.Atoi(path[1])
			if err != nil {
				next = make(map[string]interface{})
			} else {
				next = make([]interface{}, 0)
			}
		}
		source[path[0]] = mt.set(next, path[1:], v)
		return source
	case []interface{}:
		index, err := strconv.Atoi(path[0])
		if err != nil {
			return source
		}
		for len(source) <= index {
			source = append(source, nil)
		}
		next := source[index]
		if next == nil {
			_, err := strconv.Atoi(path[1])
			if err != nil {
				next = make(map[string]interface{})
			} else {
				next = make([]interface{}, 0)
			}
		}
		source[index] = mt.set(next, path[1:], v)
		return source
	default:
		return source
	}
}

func (mt *MapTree) Remove(path []string) {
	mt.remove(mt.data, path)
}

func (mt *MapTree) remove(source interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	if len(path) == 1 {
		switch source := source.(type) {
		case map[string]interface{}:
			source[path[0]] = nil
			for k, v := range source {
				if v == nil {
					delete(source, k)
				}
			}
			if len(source) == 0 {
				return nil
			}
			return source
		case []interface{}:
			index, err := strconv.Atoi(path[0])
			if err != nil || len(source) <= index {
				return source
			}
			source[index] = nil
			for i := len(source) - 1; i >= 0; i-- {
				if source[i] == nil {
					source = source[:len(source)-1]
				} else {
					break
				}
			}
			if len(source) == 0 {
				return nil
			}
			return source
		default:
			return source
		}
	}

	switch source := source.(type) {
	case map[string]interface{}:
		next, ok := source[path[0]]
		if !ok {
			return source
		}
		source[path[0]] = mt.remove(next, path[1:])
		for k, v := range source {
			if v == nil {
				delete(source, k)
			}
		}
		if len(source) == 0 {
			return nil
		}
		return source
	case []interface{}:
		index, err := strconv.Atoi(path[0])
		if err != nil || len(source) <= index {
			return source
		}
		next := source[index]
		source[index] = mt.remove(next, path[1:])
		for i := len(source) - 1; i >= 0; i-- {
			if source[i] == nil {
				source = source[:len(source)-1]
			} else {
				break
			}
		}
		if len(source) == 0 {
			return nil
		}
		return source
	default:
		return source
	}
}

func (mt *MapTree) Merge(smt *MapTree) {
	mt.data = mt.merge(smt.data, mt.data).(map[string]interface{})
}

func (mt *MapTree) merge(source interface{}, target interface{}) interface{} {
	sourceType := reflect.TypeOf(source)
	targetType := reflect.TypeOf(target)
	if sourceType != targetType {
		return source
	}
	switch source := source.(type) {
	case map[string]interface{}:
		target := target.(map[string]interface{})
		for sk, sv := range source {
			dv, ok := target[sk]
			if !ok || dv == nil {
				target[sk] = sv
			} else {
				target[sk] = mt.merge(sv, target[sk])
			}
		}
		return target
	case []interface{}:
		target := target.([]interface{})
		for index := 0; index < len(source); index++ {
			if index >= len(target) {
				target = append(target, source[index])
			} else {
				target[index] = mt.merge(source[index], target[index])
			}
		}
		return target
	default:
		return source
	}
}

func (mt *MapTree) Dulplicate() *MapTree {
	mapTree := NewMapTree()
	mapTree.SetData(deepcopy.Copy(mt.data).(map[string]interface{}))
	return mapTree
}

func (mt *MapTree) MarshalHash(cs *magic.ChainStyle) map[string]interface{} {
	hashMap := make(map[string]interface{})
	mt.marshalHash(hashMap, mt.data, "", cs)
	return hashMap
}

func (mt *MapTree) marshalHash(m map[string]interface{}, source interface{}, prefix string, cs *magic.ChainStyle) {
	switch source := source.(type) {
	case map[string]interface{}:
		for k, v := range source {
			if prefix == "" {
				mt.marshalHash(m, v, k, cs)
			} else {
				// TODO: revert Standard
				mt.marshalHash(m, v, fmt.Sprintf("%s%s%s", prefix, cs.ChainSeperator, k), cs)
			}
		}
	case []interface{}:
		for i := 0; i < len(source); i++ {
			mt.marshalHash(m, source[i], fmt.Sprintf("%s%s%d", prefix, cs.ChainSeperator, i), cs)
		}
	default:
		m[prefix] = source
	}
}

func (mt *MapTree) UnmarshalHash(hashMap map[string]interface{}, cs *magic.ChainStyle) {
	for k, v := range hashMap {
		mt.Set(cs.Chain(k), v)
	}
}
