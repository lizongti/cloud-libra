package encoding

import (
	"errors"

	"github.com/golang/protobuf/proto"
)

var (
	ErrProtobufWrongValueType = errors.New("codec protobuf converts on wrong type value")
)

type Protobuf struct{}

func init() {
	registerCodec(new(Protobuf))
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
	return MakeBytes(data), nil
}

func (*Protobuf) Unmarshal(bytes Bytes, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return ErrProtobufWrongValueType
	}
	return proto.Unmarshal(bytes.Data, pb)
}
