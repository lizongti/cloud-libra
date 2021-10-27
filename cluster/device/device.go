package device

import (
	"context"
)

type DeviceType int

const (
	DeviceTypeEmpty DeviceType = iota
	DeviceTypeBus
	DeviceTypeRouter
	DeviceTypeService
	DeviceTypeHandler
)

type Device interface {
	String() string
	LinkGateway(Device)
	Process(context.Context, Route, []byte) error
	// Devices() map[string]Device TODO
}
