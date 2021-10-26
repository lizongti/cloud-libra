package encoding

import "github.com/golang/protobuf/proto"

type ProtobufEncoding struct{}

var protobufEncoding = new(ProtobufEncoding)

func Protobuf() Encoding {
	return protobufEncoding
}

func (*ProtobufEncoding) String() string {
	return "protobuf"
}

func (s *ProtobufEncoding) Marshal(v interface{}) ([]byte, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, ErrWrongValueType
	}
	return proto.Marshal(pb)
}

func (s *ProtobufEncoding) Unmarshal(data []byte, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return ErrWrongValueType
	}
	return proto.Unmarshal(data, pb)
}
