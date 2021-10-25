package packet

type Packet interface {
	GetLength() int
	SetLength(int)
	GetData() []byte
	SetData([]byte)
	String() string
}

