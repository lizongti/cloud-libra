package component

import (
	"github.com/aceaura/libra/core/encoding"
	"github.com/aceaura/libra/core/scheduler"
)

type Service interface {
	Component

	Dispatcher() scheduler.Dispatcher
}

type ServiceBase struct {
	*ComponentBase
}

func (*ServiceBase) Dispatcher() scheduler.Dispatcher {
	return scheduler.Empty()
}

func (*ServiceBase) Encoding() encoding.Encoding {
	return encoding.Empty()
}
