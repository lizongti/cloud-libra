package cluster

import "context"

type Handler interface {
	Serve(context.Context, []byte) ([]byte, error)
}

