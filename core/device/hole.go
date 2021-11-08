package device

import (
	"context"

	"github.com/aceaura/libra/core/route"
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

func (b *hole) Process(context.Context, route.Route, []byte) error {
	return nil
}
