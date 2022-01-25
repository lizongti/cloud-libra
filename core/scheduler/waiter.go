package scheduler

type Waiter struct {
}

func NewWaiter(opt ...ApplyWaiterOption) *Waiter {
	return &Waiter{}
}

type waiterOptions struct {
}

var defaultWaiterOptions = waiterOptions{}

type ApplyWaiterOption interface {
	apply(*waiterOptions)
}

type funcWaiterOption func(*waiterOptions)

func (f funcWaiterOption) apply(opt *waiterOptions) {
	f(opt)
}

type waiterOption int

var WaiterOption waiterOption
