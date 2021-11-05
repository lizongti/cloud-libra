package device

import (
	"context"

	"github.com/aceaura/libra/magic"
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

func Name(device Device) string {
	name := device.String()
	if name == "" {
		name = magic.TypeName(device)
	}
	return name
}

func Tree(device Device) string {
	var add = func(d Device, t gotree.Tree) {}
	add = func(d Device, t gotree.Tree) {
		for _, e := range d.Extensions() {
			for _, extension := range e {
				add(extension, t.Add(Name(extension)))
			}
		}
	}

	tree := gotree.New(device.String())
	add(device, tree)
	return tree.Print()
}
