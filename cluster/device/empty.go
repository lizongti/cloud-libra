package device

import "context"

type emptyDevice int

var (
	empty = new(emptyDevice)
)

func (*emptyDevice) String() string {
	return ""
}

func (*emptyDevice) LinkGateway(Device) {}

func (*emptyDevice) Process(_ context.Context, route Route, _ []byte) error {
	return route.Error(ErrRouteDeadEnd)
}

func Empty() Device {
	return empty
}
