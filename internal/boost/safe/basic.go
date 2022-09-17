package safe

import (
	"sync/atomic"
)

type Bool struct {
	value int32
}

func (b Bool) True() {
	atomic.StoreInt32(&b.value, 1)
}

func (b Bool) False() {
	atomic.StoreInt32(&b.value, 0)
}

func (b Bool) Set(t bool) {
	if t {
		atomic.StoreInt32(&b.value, 1)
	} else {
		atomic.StoreInt32(&b.value, 0)
	}
}

func (b Bool) Get() bool {
	return atomic.LoadInt32(&b.value) > 0
}

type Int64 struct {
	value int64
}

func (i Int64) Increase() {
	atomic.AddInt64(&i.value, 1)
}

func (i Int64) Decrease() {
	atomic.AddInt64(&i.value, -1)
}

func (i Int64) Add(n int64) int64 {
	return atomic.AddInt64(&i.value, n)
}

func (i Int64) CAS(e int64, n int64) bool {
	return atomic.CompareAndSwapInt64(&i.value, e, n)
}

func (i Int64) Set(n int64) {
	atomic.StoreInt64(&i.value, n)
}

func (i Int64) Get() int64 {
	return atomic.LoadInt64(&i.value)
}
