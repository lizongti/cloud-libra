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

type Bytes struct {
	Data []byte
}

var nilBytes = MakeBytes(nil)

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

type Encoding interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type CodecSet map[string]Encoding

func newCodecSet() CodecSet {
	return make(map[string]Encoding)
}

var codecSet = newCodecSet()

func register(codecs ...Encoding) {
	codecSet.register(codecs...)
}

func (cs CodecSet) register(codecs ...Encoding) {
	for _, codec := range codecs {
		cs[magic.TypeName(codec)] = codec
	}
}

func getCodec(name string) (Encoding, error) {
	return codecSet.getCodec(name)
}

func (es CodecSet) getCodec(name string) (Encoding, error) {
	if e, ok := es[name]; ok {
		return e, nil
	}
	return nil, ErrEncodingMissingCodec
}

type Chain struct {
	encoder []string
	decoder []string
}

func NewChain(opts ...chainOpt) *Chain {
	e := &Chain{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

var nilChain = NewChain()

func Nil() *Chain {
	return nilChain
}

func (e Chain) Reverse() Chain {
	re := Chain{
		encoder: make([]string, len(e.decoder)),
		decoder: make([]string, len(e.encoder)),
	}
	lenDecoder := len(e.decoder)
	for index := 0; index < lenDecoder; index++ {
		re.encoder[index] = e.decoder[lenDecoder-1-index]
	}
	lenEncoder := len(e.encoder)
	for index := 0; index < lenEncoder; index++ {
		re.decoder[index] = e.encoder[lenEncoder-1-index]
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

func (e Chain) Marshal(v interface{}) ([]byte, error) {
	var data []byte
	for index, name := range e.encoder {
		codec, err := getCodec(name)
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

func (e Chain) Unmarshal(data []byte, v interface{}) error {
	bytes := MakeBytes(nil)
	for index, name := range e.decoder {
		codec, err := getCodec(name)
		if err != nil {
			return err
		}
		if index < len(e.decoder)-1 {
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

type chainOpt func(*Chain)
type chainOption struct{}

var ChainOption chainOption

func (chainOption) WithEncoder(path string, codecSep magic.SeparatorType, wordSep magic.SeparatorType) chainOpt {
	return func(e *Chain) {
		e.WithEncoder(path, codecSep, wordSep)
	}
}

func (e *Chain) WithEncoder(path string, codecSep magic.SeparatorType, wordSep magic.SeparatorType) *Chain {
	names := strings.Split(path, codecSep)
	for _, name := range names {
		e.encoder = append(e.encoder, magic.Standardize(name, wordSep))
	}
	return e
}

func (chainOption) WithDecoder(path string, codecSep magic.SeparatorType, wordSep magic.SeparatorType) chainOpt {
	return func(e *Chain) {
		e.WithDecoder(path, codecSep, wordSep)
	}
}

func (e *Chain) WithDecoder(path string, codecSep magic.SeparatorType, wordSep magic.SeparatorType) *Chain {
	names := strings.Split(path, codecSep)
	for _, name := range names {
		e.decoder = append(e.decoder, magic.Standardize(name, wordSep))
	}
	return e
}
