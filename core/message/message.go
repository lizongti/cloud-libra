package message

import (
	"fmt"

	"github.com/aceaura/libra/core/encoding"
	"github.com/aceaura/libra/core/route"
)

type MessageStateType int

const (
	MessageStateAssembling MessageStateType = iota
	MessageStateDispatching
)

var messageStateName = map[MessageStateType]string{
	MessageStateAssembling:  "assembling",
	MessageStateDispatching: "dispatcing",
}

func (m MessageStateType) String() string {
	if s, ok := messageStateName[m]; ok {
		return s
	}
	return fmt.Sprintf("messageStateName=%d?", int(m))
}

type Message struct {
	id       uint64
	route    route.Route
	encoding encoding.Encoding
	data     []byte
}

func NewMessage(id uint64, route route.Route, encoding encoding.Encoding, data []byte) *Message {
	return &Message{
		id:       id,
		route:    route,
		encoding: encoding,
		data:     data,
	}
}

func (m *Message) HasEncoding() bool {
	return m.encoding != encoding.Empty()
}

func (m *Message) SetEncoding(e encoding.Encoding) {
	m.encoding = e
}

func (m *Message) String() string {
	return ""
}

func (m *Message) Marshal(v interface{}) ([]byte, error) {
	return m.encoding.Marshal(v)
}

func (m *Message) Unmarshal(data []byte, v interface{}) error {
	return m.encoding.Unmarshal(data, v)
}

func (m *Message) RouteString() string {
	return m.route.String()
}

func (m *Message) Data() []byte {
	return m.data
}

func (m *Message) Forward() *Message {
	m.route = m.route.Forward()
	return m
}

func (m *Message) Position() string {
	return m.route.Name()
}

func (m *Message) State() MessageStateType {
	if m.route.Assembling() {
		return MessageStateAssembling
	}
	return MessageStateDispatching
}

func (m *Message) Reply(data []byte) *Message {
	return &Message{
		id:       m.id,
		route:    m.route.Reverse(),
		encoding: m.encoding.Reverse(),
		data:     data,
	}
}
