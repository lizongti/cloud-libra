package codec

import (
	"errors"

	"github.com/golang/protobuf/proto"
)

var (
	ErrProtobufWrongValueType = errors.New("codec protobuf converts on wrong type value")
)

type Protobuf struct{}

func init() {
	Register(new(Protobuf))
}

func (*Protobuf) String() string {
	return "protobuf"
}

func (*Protobuf) Marshal(v interface{}) (Bytes, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nilBytes, ErrProtobufWrongValueType
	}
	data, err := proto.Marshal(pb)
	if err != nil {
		return nilBytes, err
	}
	return Bytes{data}, nil
}

func (*Protobuf) Unmarshal(bytes Bytes, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return ErrProtobufWrongValueType
	}
	return proto.Unmarshal(bytes.Data, pb)
}
