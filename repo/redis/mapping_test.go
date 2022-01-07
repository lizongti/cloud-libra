package redis_test

import (
	"fmt"
	"testing"

	"github.com/aceaura/libra/core/device"
	"github.com/aceaura/libra/repo/redis"
	"github.com/alicebob/miniredis/v2"
)

func TestMapping(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	m := redis.NewMapping(
		redis.MappingOption.URL(fmt.Sprintf("redis://%s/0", s.Addr())),
		redis.MappingOption.Name("RedisMapping"),
	)
	device.Bus().WithDevice(m)

	if err := m.Replace(redis.String{
		Key:   "test_string",
		Value: "1111",
	}); err != nil {
		t.Fatal(err)
	}

	if err := m.Replace(redis.Hash{
		Key: "test_hash",
		Value: map[interface{}]interface{}{
			"a.b.c": 1,
			"b.c.d": "bcd",
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := m.Replace(redis.Set{
		Key: "test_key",
		Value: []interface{}{
			1, "2", 1, "4",
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := m.Replace(redis.Set{
		Key: "test_key",
		Value: []interface{}{
			1, "2", 1, "4",
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := m.Replace(redis.SortedSet{
		Key: "test_key",
		Value: map[interface{}]interface{}{
			"Jack":   1,
			"Jane":   1,
			"Jessie": 1,
		},
	}); err != nil {
		t.Fatal(err)
	}
}
