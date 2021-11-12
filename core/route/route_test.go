package route_test

import (
	"fmt"
	"testing"

	"github.com/aceaura/libra/core/route"
	"github.com/aceaura/libra/magic"
)

func TestRoute(t *testing.T) {
	src := magic.ChainPeriodUnderscore("bus.test_route.src")
	dst := magic.ChainPeriodUnderscore("bus.test_route.dst")
	route := route.NewRoute(src, dst)
	t.Logf("%v", route)
	if fmt.Sprint(route) != "[Bus:TestRoute:Src] -> [<Bus>:TestRoute:Dst]" {
		t.Fatal("route string is not as expected")
	}
}
