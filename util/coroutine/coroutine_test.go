package coroutine_test

import (
	"testing"

	"github.com/aceaura/libra/util/coroutine"

	"strings"
)

func interfaceSliceToStringSlice(data []interface{}) []string {
	strs := make([]string, 0)
	for _, i := range data {
		strs = append(strs, i.(string))
	}
	return strs
}

func TestCreate(t *testing.T) {
	c := coroutine.Create(func(c *coroutine.Coroutine, args ...interface{}) error {
		strs := interfaceSliceToStringSlice(args)
		out := strings.Join(strs, " ") // coroutine resume 1
		if out != "ID resume 1" {
			t.Error("ID flow error, should be ID resume 1")
		}
		t.Log(out)

		inData := c.Yield("ID", "yield", "2")
		strs = interfaceSliceToStringSlice(inData)
		out = strings.Join(strs, " ") // coroutine resume 3
		if out != "ID resume 3" {
			t.Error("ID flow error, should be ID resume 3")
		}
		t.Log(out)

		_ = c.Yield("ID", "yield", "4")
		return nil
	})

	_, ok := coroutine.TryResume(c.ID(), "ID", "resume", "0")
	if !ok {
		t.Log("Try resume test result Correct")
	}

	_, ok = coroutine.Resume(c.ID(), "ID", "resume", "1")
	if !ok {
		t.Log("Dead ID")
	}

	outData, ok := coroutine.Resume(c.ID(), "ID", "resume", "3")
	if !ok {
		t.Log("Dead ID")
	}
	strs := interfaceSliceToStringSlice(outData)
	out := strings.Join(strs, " ") // coroutine yield 2
	if out != "ID yield 2" {
		t.Error("ID flow error, should be ID yield 2")
	}
	t.Log(out)
}

func TestStart(t *testing.T) {
	c := coroutine.Wrap(func(c *coroutine.Coroutine, args ...interface{}) error {
		strs := interfaceSliceToStringSlice(args)
		out := strings.Join(strs, " ") // ID call 1
		if out != "ID call 1" {
			t.Error("ID flow error, should be ID call 1")
		}
		t.Log(out)

		inData := coroutine.Yield(c.ID(), "ID", "yield", "2")
		strs = interfaceSliceToStringSlice(inData)
		out = strings.Join(strs, " ") // ID resume 3
		if out != "ID resume 3" {
			t.Error("ID flow error, should be ID resume 3")
		}
		t.Log(out)

		_ = coroutine.Yield(c.ID(), "ID", "yield", "4")
		return nil
	})

	go func() {
		if err := coroutine.Call(c.ID(), "ID", "call", "1"); err != nil {
			t.Error(err)
		}
	}()

	outData, ok := coroutine.Resume(c.ID(), "ID", "resume", "3")
	if !ok {
		t.Log("Dead ID")
	}
	strs := interfaceSliceToStringSlice(outData)
	out := strings.Join(strs, " ") // ID yield 2
	if out != "ID yield 2" {
		t.Error("ID flow error, should be ID yield 2")
	}
	t.Log(out)

	outData, ok = coroutine.Resume(c.ID(), "ID", "resume", "5")
	if !ok {
		t.Log("Dead ID")
	}
	strs = interfaceSliceToStringSlice(outData)
	out = strings.Join(strs, " ") // ID yield 4
	if out != "ID yield 4" {
		t.Error("ID flow error, should be ID yield 4")
	}
	t.Log(out)
}
