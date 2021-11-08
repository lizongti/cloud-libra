package scheduler

import (
	"context"

	"github.com/aceaura/libra/core/route"
)

type Dispatcher interface {
	Dispatch(context.Context, route.Route) *Scheduler
}

type defaultDispatcher int

func DefaultDispatcher() Dispatcher {
	return new(defaultDispatcher)
}

func (defaultDispatcher) Dispatch(context.Context, route.Route) *Scheduler {
	return Empty()
}
