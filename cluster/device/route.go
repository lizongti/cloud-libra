package device

import (
	"errors"
	"strings"

	"github.com/aceaura/libra/magic"
)

var (
	ErrRouteDeadEnd       = errors.New("route goes to a dead end")
	ErrRouteMissingDevice = errors.New("route goes on a missing device")
)

type Route struct {
	src      []string
	dst      []string
	dstIndex int
}

func NewRoute(opts ...routeOpt) *Route {
	r := &Route{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r Route) String() string {
	var builder strings.Builder
	builder.WriteString(magic.SeparatorBracketleft)
	for index, name := range r.dst {
		builder.WriteString(name)
		if index != len(r.src)-1 {
			builder.WriteString(magic.SeparatorColon)
		}
	}
	builder.WriteString(magic.SeparatorBracketright)
	builder.WriteString(magic.SeparatorSpace)
	builder.WriteString(magic.SeparatorMinus)
	builder.WriteString(magic.SeparatorGreater)
	builder.WriteString(magic.SeparatorSpace)
	builder.WriteString(magic.SeparatorBracketleft)
	for index, name := range r.dst {
		if index == r.dstIndex {
			builder.WriteString(magic.SeparatorSpace)
			builder.WriteString(magic.SeparatorGreater)
			builder.WriteString(magic.SeparatorGreater)
			builder.WriteString(magic.SeparatorGreater)
			builder.WriteString(name)
			builder.WriteString(magic.SeparatorLess)
			builder.WriteString(magic.SeparatorLess)
			builder.WriteString(magic.SeparatorLess)
			builder.WriteString(magic.SeparatorSpace)
		} else {
			builder.WriteString(name)
		}
		if index != len(r.dst)-1 {
			builder.WriteString(magic.SeparatorColon)
		}
	}
	builder.WriteString(magic.SeparatorBracketright)
	return builder.String()
}

func (r Route) deviceType() DeviceType {
	switch r.dstIndex {
	case 0:
		return DeviceTypeBus
	case len(r.dst) - 2:
		return DeviceTypeService
	case len(r.dst) - 1:
		return DeviceTypeHandler
	default:
		return DeviceTypeRouter
	}
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

func standardize(s string, sep magic.SeparatorType) string {
	if sep == magic.SeparatorNone {
		b := []byte(s)
		if b[0] >= 'a' && b[0] <= 'z' {
			b[0] -= 32
		}
		return string(b)
	}

	b := []byte{}
	words := strings.Split(s, sep)
	for _, word := range words {
		word = standardize(word, magic.SeparatorNone)
		b = append(b, []byte(word)...)
	}
	return string(b)
}

type routeOpt func(*Route)
type routeOption struct{}

var RouteOption routeOption

func (routeOption) WithSrc(path string, deviceSep magic.SeparatorType, wordSep magic.SeparatorType) routeOpt {
	return func(r *Route) {
		r.WithSrc(path, deviceSep, wordSep)
	}
}

func (r *Route) WithSrc(path string, deviceSep magic.SeparatorType, wordSep magic.SeparatorType) *Route {
	names := strings.Split(path, deviceSep)
	for _, name := range names {
		r.src = append(r.src, standardize(name, wordSep))
	}
	return r
}

func (routeOption) WithDst(path string, deviceSep magic.SeparatorType, wordSep magic.SeparatorType) routeOpt {
	return func(r *Route) {
		r.WithDst(path, deviceSep, wordSep)
	}
}

func (r *Route) WithDst(path string, deviceSep magic.SeparatorType, wordSep magic.SeparatorType) *Route {
	names := strings.Split(path, deviceSep)
	for _, name := range names {
		r.dst = append(r.dst, standardize(name, wordSep))
	}
	return r
}

func (routeOption) WithDstIndex(dstIndex int) routeOpt {
	return func(r *Route) {
		r.WithDstIndex(dstIndex)
	}
}

func (r *Route) WithDstIndex(index int) *Route {
	r.dstIndex = index
	return r
}
