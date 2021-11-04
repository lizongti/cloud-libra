package encoding

import (
	"errors"

	"github.com/aceaura/libra/magic"
	"github.com/golang/protobuf/proto"
)

var (
	ErrProtobufWrongValueType = errors.New("codec protobuf converts on wrong type value")
)

type Protobuf struct{}

func init() {
	register(NewProtobuf())
}

func NewProtobuf() *Protobuf {
	return new(Protobuf)
}

func (protobuf Protobuf) String() string {
	return magic.TypeName(protobuf)
}

func (Protobuf) Marshal(v interface{}) ([]byte, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, ErrProtobufWrongValueType
	}
	data, err := proto.Marshal(pb)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (Protobuf) Unmarshal(data []byte, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return ErrProtobufWrongValueType
	}
	return proto.Unmarshal(data, pb)
}
