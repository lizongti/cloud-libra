package router

import (
	"context"

	"github.com/lizongti/libra/module/core/codec"
	"github.com/lizongti/libra/module/core/component"
	"github.com/lizongti/libra/module/core/handler"
)

type Router interface {
	component.Component

	Handle(string, handler.Handler)
	Serve(context.Context, string, []byte) ([]byte, error)
	WithCodec(codec.Codec)
}
