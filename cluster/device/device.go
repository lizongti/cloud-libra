package device

import (
	"context"
)

type Device interface {
	String() string
	Tree() string
	Access(Device)
	Extend(Device)
	Route(name string) Device
	Gateway() Device
	Process(context.Context, Route, []byte) error
}
