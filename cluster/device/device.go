package device

import (
	"context"
)

type DeviceType int

const (
	DeviceTypeBus DeviceType = iota
	DeviceTypeRouter
	DeviceTypeService
	DeviceTypeHandler
)

type Device interface {
	// component.Component
	String() string
	Gateway(Device)
	Process(context.Context, Route, []byte) error
}

type emptyDevice int

var (
	empty = new(emptyDevice)
)

func (*emptyDevice) String() string {
	return ""
}

func (*emptyDevice) Process(context.Context, Route, []byte) error {
	return ErrRouteDeadEnd
}

func (*emptyDevice) Gateway(Device) {}

func Empty() Device {
	return empty
}
