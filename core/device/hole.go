package device

import (
	"context"

	"github.com/aceaura/libra/core/message"
)

type hole struct {
	*Base
}

func Hole() Device {
	return newHole()
}

func newHole() *hole {
	return &hole{
		Base: NewBase(),
	}
}

func (b *hole) Process(context.Context, *message.Message) error {
	return nil
}
