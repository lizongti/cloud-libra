package redis_test

import (
	"testing"

	"github.com/aceaura/libra/repo/redis"
	"github.com/alicebob/miniredis/v2"
	"github.com/spf13/cast"
)

func TestClient(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	c := redis.NewRedis().WithAddr(s.Addr()).Client()
	originValue := "value"
	_, err = c.Do("SET", "test", originValue)
	if err != nil {
		t.Fatal(err)
	}
	value, err := c.Do("GET", "test")
	if err != nil {
		t.Fatal(err)
	}
	if string(value.([]uint8)) != originValue {
		t.Errorf("Expect redis get result equals to %s", originValue)
	}
	_, err = c.Do("SET", "test2", 1)
	if err != nil {
		t.Fatal(err)
	}
	value, err = c.Do("GET", "test2")
	if err != nil {
		t.Fatal(err)
	}
	if cast.ToInt64(string(value.([]uint8))) != 1 {
		t.Errorf("Expect redis get result equals to %s", originValue)
	}
}
