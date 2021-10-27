package device

import (
	"context"
)

type RouteType int

const (
	RouteTypeBus RouteType = iota
	RouteTypeDispatch
)

type Device interface {
	String() string
	LinkGateway(Device)
	Process(context.Context, Route, []byte) error
	// Devices() map[string]Device TODO
}
