package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/aceaura/libra/repo/redis"
	"github.com/alicebob/miniredis/v2"
)

func TestClient(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	c := redis.NewRedis().WithAddr(s.Addr()).Client()
	ctx := context.Background()
	originValue := "value"
	_, err = c.Set(ctx, "test", originValue, time.Duration(-1)).Result()
	if err != nil {
		t.Fatal(err)
	}
	value, err := c.Get(ctx, "test").Result()
	if err != nil {
		t.Fatal(err)
	}
	if value != originValue {
		t.Errorf("Expect redis get result equals to %s", originValue)
	}
}
