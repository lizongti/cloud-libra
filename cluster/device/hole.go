package device

import "context"

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

func (b *hole) Process(context.Context, Route, []byte) error {
	return nil
}
