package device

import (
	"context"

	"github.com/aceaura/libra/core/message"
)

type hole struct {
	*Base
}

func Hole() Device {
	return NewHole()
}

func NewHole() *hole {
	return &hole{
		Base: NewBase(),
	}
}

func (*hole) Process(context.Context, *message.Message) error {
	return nil
}
