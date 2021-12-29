// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package cast provides easy and safe casting in Go.
package cast

import "time"

// Bool casts an interface to a bool type.
func Bool(i interface{}) bool {
	v, _ := BoolE(i)
	return v
}

// Time casts an interface to a time.Time type.
func Time(i interface{}) time.Time {
	v, _ := TimeE(i)
	return v
}

func TimeInDefaultLocation(i interface{}, location *time.Location) time.Time {
	v, _ := TimeInDefaultLocationE(i, location)
	return v
}

// Duration casts an interface to a time.Duration type.
func Duration(i interface{}) time.Duration {
	v, _ := DurationE(i)
	return v
}

// Float64 casts an interface to a float64 type.
func Float64(i interface{}) float64 {
	v, _ := Float64E(i)
	return v
}

// Float32 casts an interface to a float32 type.
func Float32(i interface{}) float32 {
	v, _ := Float32E(i)
	return v
}

// Int64 casts an interface to an int64 type.
func Int64(i interface{}) int64 {
	v, _ := Int64E(i)
	return v
}

// Int32 casts an interface to an int32 type.
func Int32(i interface{}) int32 {
	v, _ := Int32E(i)
	return v
}

// Int16 casts an interface to an int16 type.
func Int16(i interface{}) int16 {
	v, _ := Int16E(i)
	return v
}

// Int8 casts an interface to an int8 type.
func Int8(i interface{}) int8 {
	v, _ := Int8E(i)
	return v
}

// Int casts an interface to an int type.
func Int(i interface{}) int {
	v, _ := IntE(i)
	return v
}

// Uint casts an interface to a uint type.
func Uint(i interface{}) uint {
	v, _ := UintE(i)
	return v
}

// Uint64 casts an interface to a uint64 type.
func Uint64(i interface{}) uint64 {
	v, _ := Uint64E(i)
	return v
}

// Uint32 casts an interface to a uint32 type.
func Uint32(i interface{}) uint32 {
	v, _ := Uint32E(i)
	return v
}

// Uint16 casts an interface to a uint16 type.
func Uint16(i interface{}) uint16 {
	v, _ := Uint16E(i)
	return v
}

// Uint8 casts an interface to a uint8 type.
func Uint8(i interface{}) uint8 {
	v, _ := Uint8E(i)
	return v
}

// String casts an interface to a string type.
func String(i interface{}) string {
	v, _ := StringE(i)
	return v
}

// StringMapString casts an interface to a map[string]string type.
func StringMapString(i interface{}) map[string]string {
	v, _ := StringMapStringE(i)
	return v
}

// StringMapStringSlice casts an interface to a map[string][]string type.
func StringMapStringSlice(i interface{}) map[string][]string {
	v, _ := StringMapStringSliceE(i)
	return v
}

// StringMapBool casts an interface to a map[string]bool type.
func StringMapBool(i interface{}) map[string]bool {
	v, _ := StringMapBoolE(i)
	return v
}

// StringMapInt casts an interface to a map[string]int type.
func StringMapInt(i interface{}) map[string]int {
	v, _ := StringMapIntE(i)
	return v
}

// StringMapInt64 casts an interface to a map[string]int64 type.
func StringMapInt64(i interface{}) map[string]int64 {
	v, _ := StringMapInt64E(i)
	return v
}

// StringMap casts an interface to a map[string]interface{} type.
func StringMap(i interface{}) map[string]interface{} {
	v, _ := StringMapE(i)
	return v
}

// Slice casts an interface to a []interface{} type.
func Slice(i interface{}) []interface{} {
	v, _ := SliceE(i)
	return v
}

// BoolSlice casts an interface to a []bool type.
func BoolSlice(i interface{}) []bool {
	v, _ := BoolSliceE(i)
	return v
}

// StringSlice casts an interface to a []string type.
func StringSlice(i interface{}) []string {
	v, _ := StringSliceE(i)
	return v
}

// IntSlice casts an interface to a []int type.
func IntSlice(i interface{}) []int {
	v, _ := IntSliceE(i)
	return v
}

// DurationSlice casts an interface to a []time.Duration type.
func DurationSlice(i interface{}) []time.Duration {
	v, _ := DurationSliceE(i)
	return v
}
