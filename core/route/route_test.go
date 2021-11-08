package route_test

import (
	"fmt"
	"testing"

	"github.com/aceaura/libra/cluster/route"
	"github.com/aceaura/libra/magic"
)

func TestRoute(t *testing.T) {
	route := route.NewRoute(
		route.RouteOption.WithSrc("bus.test_route.src", magic.SeparatorPeriod, magic.SeparatorUnderscore),
		route.RouteOption.WithDst("bus.test_route.dst", magic.SeparatorPeriod, magic.SeparatorUnderscore),
		route.RouteOption.WithDstIndex(1),
	)
	t.Logf("%v", route)
	if fmt.Sprint(route) != "[Bus:TestRoute:Src] -> [Bus:<TestRoute>:Dst]" {
		t.Fatal("route string is not as expected")
	}
}
