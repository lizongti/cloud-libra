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
	var result []string
	result, err = c.Command("SET", "test", originValue)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
	result, err = c.Command("GET", "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
	if cast.ToString(result[0]) != originValue {
		t.Errorf("Expect redis get result equals to %s", originValue)
	}
	result, err = c.Command("SET", "test2", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
	result, err = c.Command("GET", "test2")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
	if cast.ToInt(result[0]) != 1 {
		t.Errorf("Expect redis get result equals to %s", originValue)
	}
}
