package device

import "errors"

var (
	ErrRouteDeadEnd       = errors.New("route goes to a dead end")
	ErrRouteMissingDevice = errors.New("route goes on a missing device")
)

type Route struct {
	src      []string
	dst      []string
	dstIndex int
}

func (r Route) deviceType() DeviceType {
	if r.dstIndex == len(r.dst)-1 {
		return DeviceTypeHandler
	}
	if r.dstIndex == len(r.dst)-2 {
		return DeviceTypeService
	}
	return DeviceTypeRouter
}

func (r Route) deviceName() string {
	return r.dst[r.dstIndex]
}

func (r Route) forward() Route {
	if r.dstIndex < len(r.src)-1 {
		r.dstIndex++
	}
	return r
}

func (r Route) reverse() Route {
	return Route{
		src:      r.dst,
		dst:      r.src,
		dstIndex: 0,
	}
}
