package coroutine_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/cloudlibraries/libra/internal/boost/coroutine"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	c, _ := coroutine.Create(ctx, func(c *coroutine.Coroutine, args ...interface{}) error {
		if !reflect.DeepEqual(args, []interface{}{"Resume", "1"}) {
			t.Error("flow error, expected `Resume 1`")
		}

		in, err := c.Yield("Yield", "2")
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(in, []interface{}{"Resume", "3"}) {
			t.Error("flow error, expected `Resume 3`")
		}

		return nil
	})

	_, err := c.Resume("Resume", "1")
	if err != nil {
		t.Fatal(err)
	}

	out, err := c.Resume("Resume", "3")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(out, []interface{}{"Yield", "2"}) {
		t.Error("flow error, expected `Yield 2`")
	}
}

func TestStart(t *testing.T) {
	ctx := context.Background()
	c := coroutine.Wrap(ctx, func(c *coroutine.Coroutine, args ...interface{}) error {
		if !reflect.DeepEqual(args, []interface{}{"Call", "1"}) {
			t.Error("flow error, expected `Call 1`")
		}

		in, err := coroutine.Yield(c.ID(), "Yield", "2")
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(in, []interface{}{"Resume", "3"}) {
			t.Fatal("ID flow error, expected `Resume 3`")
		}

		_, err = coroutine.Yield(c.ID(), "Yield", "4")
		if err != nil {
			return err
		}
		return nil
	})

	go func() {
		if err := coroutine.Call(c.ID(), "Call", "1"); err != nil {
			t.Error(err)
		}
	}()

	out, err := coroutine.Resume(c.ID(), "Resume", "3")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(out, []interface{}{"Yield", "2"}) {
		t.Fatal("flow error, expected `Yield 2`")
	}

	out, err = coroutine.Resume(c.ID(), "Resume", "5")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(out, []interface{}{"Yield", "4"}) {
		t.Fatal("flow error, expected `Yield 4`")
	}
}
