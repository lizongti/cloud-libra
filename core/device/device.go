package device

import (
	"context"
	"errors"

	"github.com/aceaura/libra/core/message"
	"github.com/disiqueira/gotree"
)

var (
	ErrRouteDeadEnd       = errors.New("route has gone to a dead end")
	ErrRouteMissingDevice = errors.New("route has gone to a missing device")
)

type Device interface {
	String() string
	Access(Device)
	Extend(Device)
	Locate(name string) Device
	Gateway() Device
	Devices() map[string][]Device
	Process(context.Context, *message.Message) error
}

func Tree(device Device) string {
	var add = func(d Device, t gotree.Tree) {}
	add = func(d Device, t gotree.Tree) {
		for _, deviceGroup := range d.Devices() {
			for _, device := range deviceGroup {
				add(device, t.Add(device.String()))
			}
		}
	}

	tree := gotree.New(device.String())
	add(device, tree)
	return tree.Print()
}
