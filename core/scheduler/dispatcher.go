package scheduler

import (
	"context"

	"github.com/aceaura/libra/core/message"
)

type Dispatcher interface {
	Dispatch(context.Context, *message.Message) *Scheduler
}

type defaultDispatcher int

func DefaultDispatcher() Dispatcher {
	return new(defaultDispatcher)
}

func (defaultDispatcher) Dispatch(context.Context, *message.Message) *Scheduler {
	return Empty()
}
