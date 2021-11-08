package message

import (
	"github.com/aceaura/libra/core/encoding"
	"github.com/aceaura/libra/core/route"
)

type Message struct {
	ID       uint64
	Route    route.Route
	Encoding encoding.Encoding
	Data     []byte
}
