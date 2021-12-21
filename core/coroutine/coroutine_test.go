package coroutine_test

import (
	"reflect"
	"testing"

	"github.com/aceaura/libra/core/coroutine"
)

func TestCreate(t *testing.T) {
	c, _ := coroutine.Create(func(c *coroutine.Coroutine, args ...interface{}) error {
		if !reflect.DeepEqual(args, []interface{}{"Resume", "1"}) {
			t.Error("flow error, expected `Resume 1`")
		}

		inData, err := c.Yield("Yield", "2")
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(inData, []interface{}{"Resume", "3"}) {
			t.Error("flow error, expected `Resume 3`")
		}

		_, err = c.Yield("Yield", "4")
		if err != nil {
			t.Fatal(err)
		}
		return nil
	})

	_, err := c.TryResume("Resume", "0")
	if err != coroutine.ErrCoroutineNotSuspended {
		t.Fatal(err)
	}

	_, err = c.Resume("Resume", "1")
	if err != nil {
		t.Fatal(err)
	}

	outData, err := c.Resume("Resume", "3")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(outData, []interface{}{"Yield", "2"}) {
		t.Error("flow error, expected `Yield 2`")
	}
}

func TestStart(t *testing.T) {
	c := coroutine.Wrap(func(c *coroutine.Coroutine, args ...interface{}) error {
		if !reflect.DeepEqual(args, []interface{}{"Call", "1"}) {
			t.Error("flow error, expected `Call 1`")
		}

		inData, err := coroutine.Yield(c.ID(), "Yield", "2")
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(inData, []interface{}{"Resume", "3"}) {
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

	outData, err := coroutine.Resume(c.ID(), "Resume", "3")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(outData, []interface{}{"Yield", "2"}) {
		t.Fatal("flow error, expected `Yield 2`")
	}

	outData, err = coroutine.Resume(c.ID(), "Resume", "5")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(outData, []interface{}{"Yield", "4"}) {
		t.Fatal("flow error, expected `Yield 4`")
	}
}
