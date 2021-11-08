package dispatcher

import (
	"context"

	"github.com/aceaura/libra/core/route"
	"github.com/aceaura/libra/core/scheduler"
)

type Dispatcher interface {
	Dispatch(context.Context, route.Route) *scheduler.Scheduler
}

type defaultDispatcher int

func Default() Dispatcher {
	return new(defaultDispatcher)
}

func (defaultDispatcher) Dispatch(context.Context, route.Route) *scheduler.Scheduler {
	return scheduler.Empty()
}
