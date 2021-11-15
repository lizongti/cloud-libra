package route_test

import (
	"fmt"
	"testing"

	"github.com/aceaura/libra/core/route"
	"github.com/aceaura/libra/magic"
)

func TestRoute(t *testing.T) {
	style := magic.ChainStyle{
		ChainSeperator: magic.SeparatorPeriod,
		WordSeparator:  magic.SeparatorUnderscore,
	}
	src := style.Chain("bus.test_route.src")
	dst := style.Chain("bus.test_route.dst")
	route := route.New(src, dst)
	t.Logf("%v", route)
	if fmt.Sprint(route) != "[Bus:TestRoute:Src] -> [<Bus>:TestRoute:Dst]" {
		t.Fatal("route string is not as expected")
	}
}
