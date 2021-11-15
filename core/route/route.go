package route

import (
	"fmt"
	"strings"

	"github.com/aceaura/libra/magic"
)

type Route struct {
	src   []string
	dst   []string
	index int
}

func New(src, dst []string) *Route {
	r := &Route{
		src: src,
		dst: dst,
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
		if index == r.index {
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
	return r.index > 0
}

func (r Route) Assembling() bool {
	return r.index == 0
}

func (r Route) Forward() Route {
	if r.index < len(r.dst)-1 {
		r.index++
	}
	return r
}

func (r Route) Position() string {
	return r.dst[r.index]
}

func (r Route) Reverse() Route {
	return Route{
		src:   r.dst,
		dst:   r.src,
		index: 0,
	}
}

func (r Route) Error(err error) error {
	return fmt.Errorf("route %v error: %w", r, err)
}
