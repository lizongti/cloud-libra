package hierarchy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var ErrHierachyShouldBeMap = errors.New("hierarchy should be map")

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

func ReadEnv(prefix string) error {
	return _default.ReadEnv(prefix)
}

func (h *Hierarchy) ReadEnv(prefix string) error {
	h.AutomaticEnv()
	h.SetEnvPrefix(prefix)

	return nil
}

func ReadFlags( ) error {
	return _default.ReadFlags()
}

func (h *Hierarchy) ReadFlags() error {

func ReadAssetMap(assetMap map[string][]byte) error {
	return _default.ReadAssetMap(assetMap)
}

func (h *Hierarchy) ReadAssetMap(assetMap map[string][]byte) error {
	keys := make([]string, 0, len(assetMap))
	for name := range assetMap {
		keys = append(keys, name)
	}

	sort.Strings(keys)

	for _, name := range keys {
		ext := filepath.Ext(name)
		data := assetMap[name]

		v := viper.New()
		v.SetConfigType(ext[1:])

		if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
			return err
		}

		if err := h.MergeConfigMap(v.AllSettings()); err != nil {
			return err
		}
	}

	return nil
}
