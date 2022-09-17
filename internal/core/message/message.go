package message

import (
	"github.com/cloudlibraries/libra/internal/core/encoding"
	"github.com/cloudlibraries/libra/internal/core/route"
)

type Message struct {
	ID       uint64
	Route    route.Route
	Encoding encoding.Encoding
	Data     []byte
}
