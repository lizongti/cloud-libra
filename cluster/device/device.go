package device

import (
	"context"

	"github.com/disiqueira/gotree"
)

type Device interface {
	String() string
	Access(Device)
	Extend(Device)
	Route(name string) Device
	Gateway() Device
	Extensions() map[string][]Device
	Process(context.Context, Route, []byte) error
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
