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

func (*emptyDevice) Process(context.Context, Route, []byte) error {
	return ErrRouteDeadEnd
}

func Empty() Device {
	return empty
}
