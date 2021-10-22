package service

import (
	"github.com/aceaura/libra/cluster/component"
)

type Router interface {
	component.Component

	// Handle(string, handler.Handler)
	// Serve(context.Context, string, []byte) ([]byte, error)
	// WithCodec(codec.Codec)
}
