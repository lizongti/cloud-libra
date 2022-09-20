package hierarchy

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var (
	ErrHierachyShouldBeMap = errors.New("hierarchy should be map")
)

type Hierarchy struct {
	*viper.Viper
}

func NewHierarcky() *Hierarchy {
	return &Hierarchy{viper.New()}
}

var _default = NewHierarcky()

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

func (h *Hierarchy) IsArray(key string) bool {
	data, err := h.JSON()
	if err != nil {
		return false
	}
	return gjson.Get(string(data), key).IsArray()
}

func (h *Hierarchy) IsMap(key string) bool {
	data, err := h.JSON()
	if err != nil {
		return false
	}
	return gjson.Get(string(data), ".").IsObject()
}

func (h *Hierarchy) ChildrenInArray(key string) ([]*Hierarchy, error) {
	data, err := h.JSON()
	if err != nil {
		return nil, err
	}
	var children []*Hierarchy
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

func (h *Hierarchy) ChildrenInMap(key string) (map[string]*Hierarchy, error) {
	data, err := h.JSON()
	if err != nil {
		return nil, err
	}

	var children = make(map[string]*Hierarchy)
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
