package hierarchy

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

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

func (h *Hierarchy) String() string {
	data, err := h.JSONIndent()
	if err != nil {
		log.Panic(fmt.Errorf("failed to marshal hierarchy: %w", err))
	}

	return string(data)
}

func Sub(key string) *Hierarchy {
	return _default.Sub(key)
}

func (h *Hierarchy) Sub(key string) *Hierarchy {
	v := h.Viper.Sub(key)
	if v == nil {
		return &Hierarchy{viper.New()}
	}
	return &Hierarchy{v}
}

func JSON() ([]byte, error) {
	return _default.JSON()
}

func (h *Hierarchy) JSON() ([]byte, error) {
	return json.Marshal(h.AllSettings())
}

func JSONIndent() ([]byte, error) {
	return _default.JSONIndent()
}

func (h *Hierarchy) JSONIndent() ([]byte, error) {
	return json.MarshalIndent(h.AllSettings(), "", "  ")
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

	return gjson.Get(string(data), key).IsObject()
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

			children = append(children, h.Sub(childKey))
		}

		return children, nil
	}

	if node.IsObject() {
		for nodeKey, child := range node.Map() {
			childKey := fmt.Sprintf("%s.%s", key, nodeKey)
			if child.IsArray() {
				return nil, fmt.Errorf("%w: %s", ErrHierachyShouldBeMap, childKey)
			}

			children = append(children, h.Sub(childKey))
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

			children[childKey] = h.Sub(childKey)
		}

		return children, nil
	}

	if node.IsObject() {
		for nodeKey, child := range node.Map() {
			childKey := fmt.Sprintf("%s.%s", key, nodeKey)
			if child.IsArray() {
				return nil, fmt.Errorf("%w: %s", ErrHierachyShouldBeMap, childKey)
			}

			children[childKey] = h.Sub(childKey)
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

// varExp gets var from ${var}.
var varExp = regexp.MustCompile(`\$\{([a-zA-Z0-9_]+)\}`)

func (h *Hierarchy) ReplaceAllVars(data []byte) []byte {
	// replace all vars
	return varExp.ReplaceAllFunc(data, func(match []byte) []byte {
		key := string(match[2 : len(match)-1])
		return []byte(h.GetString(key))
	})
}

func GetIntVal(key string, defaultValue int) int {
	return _default.GetIntVal(key, defaultValue)
}

func (h *Hierarchy) GetIntVal(key string, defaultValue int) int {
	v := h.GetInt(key)
	if v == 0 {
		return defaultValue
	}
	return v
}

func GetInt32Val(key string, defaultValue int32) int32 {
	return _default.GetInt32Val(key, defaultValue)
}

func (h *Hierarchy) GetInt32Val(key string, defaultValue int32) int32 {
	v := h.GetInt32(key)
	if v == 0 {
		return defaultValue
	}
	return v
}

func GetInt64Val(key string, defaultValue int64) int64 {
	return _default.GetInt64Val(key, defaultValue)
}

func (h *Hierarchy) GetInt64Val(key string, defaultValue int64) int64 {
	v := h.GetInt64(key)
	if v == 0 {
		return defaultValue
	}
	return v
}

func GetUintVal(key string, defaultValue uint) uint {
	return _default.GetUintVal(key, defaultValue)
}

func (h *Hierarchy) GetUintVal(key string, defaultValue uint) uint {
	v := h.GetUint(key)
	if v == 0 {
		return defaultValue
	}
	return v
}

func GetUint32Val(key string, defaultValue uint32) uint32 {
	return _default.GetUint32Val(key, defaultValue)
}

func (h *Hierarchy) GetUint32Val(key string, defaultValue uint32) uint32 {
	v := h.GetUint32(key)
	if v == 0 {
		return defaultValue
	}
	return v
}

func GetUint64Val(key string, defaultValue uint64) uint64 {
	return _default.GetUint64Val(key, defaultValue)
}

func (h *Hierarchy) GetUint64Val(key string, defaultValue uint64) uint64 {
	v := h.GetUint64(key)
	if v == 0 {
		return defaultValue
	}
	return v
}

func GetFloat64Val(key string, defaultValue float64) float64 {
	return _default.GetFloat64Val(key, defaultValue)
}

func (h *Hierarchy) GetFloat64Val(key string, defaultValue float64) float64 {
	v := h.GetFloat64(key)
	if v == 0 {
		return defaultValue
	}
	return v
}

func GetBoolVal(key string, defaultValue bool) bool {
	return _default.GetBoolVal(key, defaultValue)
}

func (h *Hierarchy) GetBoolVal(key string, defaultValue bool) bool {
	v := h.GetBool(key)
	if !v {
		return defaultValue
	}
	return v
}

func GetStringVal(key string, defaultValue string) string {
	return _default.GetStringVal(key, defaultValue)
}

func (h *Hierarchy) GetStringVal(key string, defaultValue string) string {
	v := h.GetString(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func GetDurationVal(key string, defaultValue time.Duration) time.Duration {
	return _default.GetDurationVal(key, defaultValue)
}

func (h *Hierarchy) GetDurationVal(key string, defaultValue time.Duration) time.Duration {
	v := h.GetDuration(key)
	if v == 0 {
		return defaultValue
	}
	return v
}

func GetTimeVal(key string, defaultValue time.Time) time.Time {
	return _default.GetTimeVal(key, defaultValue)
}

func (h *Hierarchy) GetTimeVal(key string, defaultValue time.Time) time.Time {
	v := h.GetTime(key)
	if v.IsZero() {
		return defaultValue
	}
	return v
}

func GetIntSliceVal(key string, defaultValue []int) []int {
	return _default.GetIntSliceVal(key, defaultValue)
}

func (h *Hierarchy) GetIntSliceVal(key string, defaultValue []int) []int {
	v := h.GetIntSlice(key)
	if len(v) == 0 {
		return defaultValue
	}
	return v
}

func GetStringMapVal(key string, defaultValue map[string]interface{}) map[string]interface{} {
	return _default.GetStringMapVal(key, defaultValue)
}

func (h *Hierarchy) GetStringMapVal(key string, defaultValue map[string]interface{}) map[string]interface{} {
	v := h.GetStringMap(key)
	if len(v) == 0 {
		return defaultValue
	}
	return v
}

func GetStringMapStringVal(key string, defaultValue map[string]string) map[string]string {
	return _default.GetStringMapStringVal(key, defaultValue)
}

func (h *Hierarchy) GetStringMapStringVal(key string, defaultValue map[string]string) map[string]string {
	v := h.GetStringMapString(key)
	if len(v) == 0 {
		return defaultValue
	}
	return v
}

func GetStringMapStringSliceVal(key string, defaultValue map[string][]string) map[string][]string {
	return _default.GetStringMapStringSliceVal(key, defaultValue)
}

func (h *Hierarchy) GetStringMapStringSliceVal(key string, defaultValue map[string][]string) map[string][]string {
	v := h.GetStringMapStringSlice(key)
	if len(v) == 0 {
		return defaultValue
	}
	return v
}

func GetSizeInBytesVal(key string, defaultValue uint) uint {
	return _default.GetSizeInBytesVal(key, defaultValue)
}

func (h *Hierarchy) GetSizeInBytesVal(key string, defaultValue uint) uint {
	v := h.GetSizeInBytes(key)
	if v == 0 {
		return defaultValue
	}
	return v
}
