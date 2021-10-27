package device_test

import (
	"fmt"
	"testing"

	"github.com/aceaura/libra/cluster/device"
	"github.com/aceaura/libra/magic"
)

func TestRoute(t *testing.T) {
	route := device.NewRoute(
		device.RouteOption.WithSrc("bus.test_route.src", magic.SeparatorPeriod, magic.SeparatorUnderscore),
		device.RouteOption.WithDst("bus.test_route.dst", magic.SeparatorPeriod, magic.SeparatorUnderscore),
		device.RouteOption.WithDstIndex(1),
	)
	t.Logf("%v", route)
	if fmt.Sprint(route) != "[Bus:TestRoute:Dst] -> [Bus: >>>TestRoute<<< :Dst]" {
		t.Fatal("route string is not as expected")
	}
}
