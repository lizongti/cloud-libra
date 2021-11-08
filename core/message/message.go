package message

import (
	"github.com/aceaura/libra/core/encoding"
	"github.com/aceaura/libra/core/route"
)

type Message struct {
	id    uint64
	route route.Route
	codec encoding.Codec
	data  []byte
}

func (m *Message) Data() []byte {
	return m.data
}

func (m *Message) Route() route.Route {
	return m.route
}

func (m *Message) ID() uint64 {
	return m.id
}

func Reply(m *Message, data []byte) *Message {
	return &Message{
		id:    m.id,
		route: m.route.Reverse(),
		codec: m.codec.Reverse(),
		data:  data,
	}
}
