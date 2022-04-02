package redis_test

import (
	"fmt"
	"testing"

	"github.com/aceaura/libra/boost/cast"
	"github.com/aceaura/libra/core/device"
	"github.com/aceaura/libra/core/encoding"
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
		fmt.Sprintf("redis://%s/0", s.Addr()),
		redis.WithMappingName("RedisMapping"),
	)
	device.Bus().Integrate(m)

	testData := []interface{}{
		redis.String{
			Key:   "test_empty_string",
			Value: "",
		},
		redis.String{
			Key:   "test_string",
			Value: "1111",
		},
		redis.Hash{
			Key:   "test_empty_hash",
			Value: map[string]string{},
		},
		redis.Hash{
			Key: "test_hash",
			Value: cast.ToStringMapString(map[string]interface{}{
				"a.b.c": 1,
				"b.c.d": "bcd",
			}),
		},
		redis.List{
			Key:   "test_empty_list",
			Value: []string{},
		},
		redis.List{
			Key: "test_list",
			Value: cast.ToStringSlice([]interface{}{
				1, "2", 3, "4",
			}),
		},
		redis.Set{
			Key:   "test_empty_set",
			Value: []string{},
		},
		redis.Set{
			Key: "test_set",
			Value: cast.ToStringSlice([]interface{}{
				1, "2", 3, "4",
			}),
		},
		redis.SortedSet{
			Key:   "test_empty_sorted_set",
			Value: map[string]string{},
		},
		redis.SortedSet{
			Key: "test_sorted_set",
			Value: cast.ToStringMapString(map[interface{}]interface{}{
				"Jack":   1,
				"Jane":   1,
				"Jessie": 1,
			}),
		},
	}

	for _, v := range testData {
		if err := m.Replace(v); err != nil {
			t.Fatal(err)
		}
		v2, err := m.Select(v)
		if err != nil {
			t.Fatal(err)
		}
		data, err := encoding.Marshal(encoding.NewJSON(), v)
		if err != nil {
			t.Fatal(err)
		}
		data2, err := encoding.Marshal(encoding.NewJSON(), v2)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != string(data2) {
			t.Fatalf("expected select result equals origin result. select: %s, origin: %s", string(data2), string(data))
		}
	}
}
