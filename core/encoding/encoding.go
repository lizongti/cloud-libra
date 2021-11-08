package encoding

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aceaura/libra/magic"
)

var (
	ErrEncodingMissingCodec = errors.New("encoding cannot find codec by name")
)

type Encoding interface {
	String() string
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type Bytes struct {
	Data []byte
}

var nilBytes = MakeBytes(nil)

func NewBytes() *Bytes {
	return new(Bytes)
}

func MakeBytes(v interface{}) Bytes {
	switch v := v.(type) {
	case []byte:
		return Bytes{
			Data: v,
		}
	case string:
		return Bytes{
			Data: []byte(v),
		}
	case Bytes:
		return v.Dulplicate()
	case *Bytes:
		return v.Dulplicate()
	default:
		return Bytes{
			Data: []byte(fmt.Sprint(v)),
		}
	}
}

func (b Bytes) Dulplicate() Bytes {
	data := make([]byte, len(b.Data))
	copy(data, b.Data)
	return MakeBytes(data)
}

func (b Bytes) Copy(in Bytes) {
	b.Data = make([]byte, len(in.Data))
	copy(b.Data, in.Data)
}

type EncodingSet map[string]Encoding

func newEncodingSet() EncodingSet {
	return make(map[string]Encoding)
}

var encodingSet = newEncodingSet()

func register(encodings ...Encoding) {
	encodingSet.register(encodings...)
}

func (es EncodingSet) register(codecs ...Encoding) {
	for _, codec := range codecs {
		es[magic.TypeName(codec)] = codec
	}
}

func codec(name string) (Encoding, error) {
	return encodingSet.codec(name)
}

func (es EncodingSet) codec(name string) (Encoding, error) {
	if e, ok := es[name]; ok {
		return e, nil
	}
	return nil, ErrEncodingMissingCodec
}

type Codec struct {
	encoder []string
	decoder []string
}

func NewCodec(opts ...funcCodecOption) *Codec {
	e := &Codec{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

var emptyCodec = NewCodec()

func Empty() Codec {
	return *emptyCodec
}

func (c Codec) String() string {
	var builder strings.Builder
	builder.WriteString(magic.SeparatorBracketleft)
	for index, name := range c.encoder {
		builder.WriteString(name)
		if index != len(c.encoder)-1 {
			builder.WriteString(magic.SeparatorColon)
		}
	}
	builder.WriteString(magic.SeparatorBracketright)
	builder.WriteString(magic.SeparatorSpace)
	builder.WriteString(magic.SeparatorMinus)
	builder.WriteString(magic.SeparatorGreater)
	builder.WriteString(magic.SeparatorSpace)
	builder.WriteString(magic.SeparatorBracketleft)
	for index, name := range c.decoder {
		builder.WriteString(name)
		if index != len(c.decoder)-1 {
			builder.WriteString(magic.SeparatorColon)
		}
	}
	builder.WriteString(magic.SeparatorBracketright)
	return builder.String()
}

func (c Codec) Reverse() Codec {
	re := Codec{
		encoder: make([]string, len(c.decoder)),
		decoder: make([]string, len(c.encoder)),
	}
	lenDecoder := len(c.decoder)
	for index := 0; index < lenDecoder; index++ {
		re.encoder[index] = c.decoder[lenDecoder-1-index]
	}
	lenEncoder := len(c.encoder)
	for index := 0; index < lenEncoder; index++ {
		re.decoder[index] = c.encoder[lenEncoder-1-index]
	}
	return re
}

func Encode(e Encoding, v interface{}) ([]byte, error) {
	return e.Marshal(v)
}

func Decode(e Encoding, data []byte, v interface{}) error {
	return e.Unmarshal(data, v)
}

func Marshal(e Encoding, v interface{}) ([]byte, error) {
	return e.Marshal(v)
}

func Unmarshal(e Encoding, data []byte, v interface{}) error {
	return e.Unmarshal(data, v)
}

func (c Codec) Marshal(v interface{}) ([]byte, error) {
	var data []byte
	for index, name := range c.encoder {
		codec, err := codec(name)
		if err != nil {
			return nil, err
		}
		if index == 0 {
			data, err = codec.Marshal(v)
			if err != nil {
				return nil, err
			}
		} else {
			data, err = codec.Marshal(data)
			if err != nil {
				return nil, err
			}
		}

	}
	return data, nil
}

func (c Codec) Unmarshal(data []byte, v interface{}) error {
	bytes := MakeBytes(nil)
	for index, name := range c.decoder {
		codec, err := codec(name)
		if err != nil {
			return err
		}
		if index < len(c.decoder)-1 {
			err = codec.Unmarshal(data, &bytes)
			if err != nil {
				return err
			}
			data = bytes.Data
		} else {
			err = codec.Unmarshal(data, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type funcCodecOption func(*Codec)
type chainOption struct{}

var ChainOption chainOption

func (chainOption) WithEncoder(path string, codecSep magic.SeparatorType, wordSep magic.SeparatorType) funcCodecOption {
	return func(c *Codec) {
		c.WithEncoder(path, codecSep, wordSep)
	}
}

func (c *Codec) WithEncoder(path string, codecSep magic.SeparatorType, wordSep magic.SeparatorType) *Codec {
	names := strings.Split(path, codecSep)
	for _, name := range names {
		c.encoder = append(c.encoder, magic.Standardize(name, wordSep))
	}
	return c
}

func (chainOption) WithDecoder(path string, codecSep magic.SeparatorType, wordSep magic.SeparatorType) funcCodecOption {
	return func(c *Codec) {
		c.WithDecoder(path, codecSep, wordSep)
	}
}

func (c *Codec) WithDecoder(path string, codecSep magic.SeparatorType, wordSep magic.SeparatorType) *Codec {
	names := strings.Split(path, codecSep)
	for _, name := range names {
		c.decoder = append(c.decoder, magic.Standardize(name, wordSep))
	}
	return c
}
