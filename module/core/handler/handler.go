package handler

import "github.com/lizongti/libra/context"

type Handler interface {
	Serve(context.Context, []byte) ([]byte, error)
	String() string
}
