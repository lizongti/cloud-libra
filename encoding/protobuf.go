package encoding

import (
	"errors"

	"github.com/golang/protobuf/proto"
)

var (
	ErrProtobufWrongValueType = errors.New("encoding convert on wrong type value")
)

type ProtobufEncoding struct{}

var protobufEncoding = new(ProtobufEncoding)

func Protobuf() TypeEncoding {
	return protobufEncoding
}

func (*ProtobufEncoding) String() string {
	return "protobuf"
}

func (*ProtobufEncoding) Marshal(v interface{}) ([]byte, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, ErrProtobufWrongValueType
	}
	return proto.Marshal(pb)
}

func (*ProtobufEncoding) Unmarshal(data []byte, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return ErrProtobufWrongValueType
	}
	return proto.Unmarshal(data, pb)
}
