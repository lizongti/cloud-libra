package device

import (
	"context"
	"errors"
	"fmt"

	"github.com/aceaura/libra/core/message"
	"github.com/disiqueira/gotree"
)

var (
	ErrRouteDeadEnd       = errors.New("route has gone to a dead end")
	ErrRouteMissingDevice = errors.New("route has gone to a missing device")
)

func routeErr(routeStr string, err error) error {
	return fmt.Errorf("route %s error: %w", routeStr, err)
}

type Device interface {
	String() string
	Access(Device)
	Extend(Device)
	Locate(name string) Device
	Gateway() Device
	Extensions() map[string][]Device
	Process(context.Context, *message.Message) error
}

func Tree(device Device) string {
	var add = func(d Device, t gotree.Tree) {}
	add = func(d Device, t gotree.Tree) {
		for _, e := range d.Extensions() {
			for _, extension := range e {
				add(extension, t.Add(extension.String()))
			}
		}
	}

	tree := gotree.New(device.String())
	add(device, tree)
	return tree.Print()
}
