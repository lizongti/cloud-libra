package hierarchy

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var ErrHierachyShouldBeMap = errors.New("hierarchy should be map")

type Hierarchy struct {
	*viper.Viper
}

func New() *Hierarchy {
	return &Hierarchy{viper.New()}
}

var _default = New()

func Child(key string) *Hierarchy {
	return _default.Child(key)
}

func (h *Hierarchy) Child(key string) *Hierarchy {
	return &Hierarchy{h.Viper.Sub(key)}
}

func JSON() ([]byte, error) {
	return _default.JSON()
}

func (h *Hierarchy) JSON() ([]byte, error) {
	return json.Marshal(h.AllSettings())
}

func IsArray(key string) bool {
	return _default.IsArray(key)
}

func (h *Hierarchy) IsArray(key string) bool {
	data, err := h.JSON()
	if err != nil {
		return false
	}

	return gjson.Get(string(data), key).IsArray()
}

func IsMap(key string) bool {
	return _default.IsMap(key)
}

func (h *Hierarchy) IsMap(key string) bool {
	data, err := h.JSON()
	if err != nil {
		return false
	}

	return gjson.Get(string(data), ".").IsObject()
}

func ChildrenInArray(key string) ([]*Hierarchy, error) {
	return _default.ChildrenInArray(key)
}

func (h *Hierarchy) ChildrenInArray(key string) ([]*Hierarchy, error) {
	data, err := h.JSON()
	if err != nil {
		return nil, err
	}

	children := make([]*Hierarchy, 0)
	node := gjson.Get(string(data), key)

	if node.IsArray() {
		for index, child := range node.Array() {
			childKey := fmt.Sprintf("%s.%d", key, index)
			if child.IsArray() {
				return nil, fmt.Errorf("%w: %s", ErrHierachyShouldBeMap, childKey)
			}

			children = append(children, h.Child(childKey))
		}

		return children, nil
	}

	if node.IsObject() {
		for nodeKey, child := range node.Map() {
			childKey := fmt.Sprintf("%s.%s", key, nodeKey)
			if child.IsArray() {
				return nil, fmt.Errorf("%w: %s", ErrHierachyShouldBeMap, childKey)
			}

			children = append(children, h.Child(childKey))
		}

		return children, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrHierachyShouldBeMap, key)
}

func ChildrenInMap(key string) (map[string]*Hierarchy, error) {
	return _default.ChildrenInMap(key)
}

func (h *Hierarchy) ChildrenInMap(key string) (map[string]*Hierarchy, error) {
	data, err := h.JSON()
	if err != nil {
		return nil, err
	}

	children := make(map[string]*Hierarchy)

	node := gjson.Get(string(data), key)
	if node.IsArray() {
		for index, child := range node.Array() {
			childKey := fmt.Sprintf("%s.%d", key, index)
			if child.IsArray() {
				return nil, fmt.Errorf("%w: %s", ErrHierachyShouldBeMap, childKey)
			}

			children[childKey] = h.Child(childKey)
		}

		return children, nil
	}

	if node.IsObject() {
		for nodeKey, child := range node.Map() {
			childKey := fmt.Sprintf("%s.%s", key, nodeKey)
			if child.IsArray() {
				return nil, fmt.Errorf("%w: %s", ErrHierachyShouldBeMap, childKey)
			}

			children[childKey] = h.Child(childKey)
		}

		return children, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrHierachyShouldBeMap, key)
}

func ForeachInArray(key string, fn func(index int, child *Hierarchy) (bool, error)) error {
	return _default.ForeachInArray(key, fn)
}

func (h *Hierarchy) ForeachInArray(key string, fn func(index int, child *Hierarchy) (bool, error)) error {
	children, err := h.ChildrenInArray(key)
	if err != nil {
		return err
	}

	for index, child := range children {
		if ok, err := fn(index, child); err != nil {
			return err
		} else if !ok {
			break
		}
	}

	return nil
}

func ForeachInMap(key string, fn func(key string, child *Hierarchy) (bool, error)) error {
	return _default.ForeachInMap(key, fn)
}

func (h *Hierarchy) ForeachInMap(key string, fn func(key string, child *Hierarchy) (bool, error)) error {
	children, err := h.ChildrenInMap(key)
	if err != nil {
		return err
	}

	for key, child := range children {
		if ok, err := fn(key, child); err != nil {
			return err
		} else if !ok {
			break
		}
	}

	return nil
}

var (
	// varExp gets var from ${var}
	varExp = regexp.MustCompile(`\$\{([a-zA-Z0-9_]+)\}`)
)

func (h *Hierarchy) ReplaceAllVars(data []byte) []byte {
	// replace all vars
	return varExp.ReplaceAllFunc(data, func(match []byte) []byte {
		key := string(match[2 : len(match)-1])
		return []byte(h.GetString(key))
	})
}
