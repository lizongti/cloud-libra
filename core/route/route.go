package route

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aceaura/libra/magic"
)

var (
	ErrRouteDeadEnd       = errors.New("route has gone to a dead end")
	ErrRouteMissingDevice = errors.New("route has gone to a missing device")
)

type Route struct {
	src      []string
	dst      []string
	dstIndex int
}

func NewRoute(opts ...funcRouteOption) *Route {
	r := &Route{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r Route) String() string {
	var builder strings.Builder
	builder.WriteString(magic.SeparatorBracketleft)
	for index, name := range r.src {
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
			builder.WriteString(magic.SeparatorLess)
			builder.WriteString(name)
			builder.WriteString(magic.SeparatorGreater)
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

func (r Route) Dispatching() bool {
	return r.dstIndex > 0
}

func (r Route) Assembling() bool {
	return r.dstIndex == 0
}

func (r Route) Forward() Route {
	if r.dstIndex < len(r.dst)-1 {
		r.dstIndex++
	}
	return r
}

func (r Route) Name() string {
	return r.dst[r.dstIndex]
}

func (r Route) Reverse() Route {
	return Route{
		src:      r.dst,
		dst:      r.src,
		dstIndex: 0,
	}
}

func (r Route) Error(err error) error {
	return fmt.Errorf("route %v error: %w", r, err)
}

type funcRouteOption func(*Route)
type routeOption struct{}

var RouteOption routeOption

func (routeOption) WithSrc(path string, deviceSep magic.SeparatorType, wordSep magic.SeparatorType) funcRouteOption {
	return func(r *Route) {
		r.WithSrc(path, deviceSep, wordSep)
	}
}

func (r *Route) WithSrc(path string, deviceSep magic.SeparatorType, wordSep magic.SeparatorType) *Route {
	names := strings.Split(path, deviceSep)
	for _, name := range names {
		r.src = append(r.src, magic.Standardize(name, wordSep))
	}
	return r
}

func (routeOption) WithDst(path string, deviceSep magic.SeparatorType, wordSep magic.SeparatorType) funcRouteOption {
	return func(r *Route) {
		r.WithDst(path, deviceSep, wordSep)
	}
}

func (r *Route) WithDst(path string, deviceSep magic.SeparatorType, wordSep magic.SeparatorType) *Route {
	names := strings.Split(path, deviceSep)
	for _, name := range names {
		r.dst = append(r.dst, magic.Standardize(name, wordSep))
	}
	return r
}

func (routeOption) WithDstIndex(dstIndex int) funcRouteOption {
	return func(r *Route) {
		r.WithDstIndex(dstIndex)
	}
}

func (r *Route) WithDstIndex(index int) *Route {
	r.dstIndex = index
	return r
}
