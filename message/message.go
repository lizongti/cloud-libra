package message

type Message interface {
	GetType() byte
	SetType(byte)
	GetID() uint64
	SetID() uint64
	GetRoute() string
	SetRoute(string)
	GetVersion() string
	SetVersion(string)
	GetData() []byte
	SetData([]byte)
	String() string
}
