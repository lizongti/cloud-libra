package scheduler

import (
	"context"

	"github.com/aceaura/libra/core/message"
)

type Dispatcher interface {
	Dispatch(context.Context, *message.Message) *Scheduler
}
