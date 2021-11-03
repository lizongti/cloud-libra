package encoding

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/aceaura/libra/magic"
)

var (
	ErrMarshallerUnknownType   = errors.New("unknown marshaller type is used")
	ErrUnmarshallerUnknownType = errors.New("unknown unmarshaller type is used")
	ErrEncodingMissingCodec    = errors.New("encoding cannot find codec by name")
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

type Codec interface {
	Marshal(interface{}) (Bytes, error)
	Unmarshal(Bytes, interface{}) error
}

type CodecSet map[string]Codec

func newCodeSet() CodecSet {
	return make(map[string]Codec)
}

var codecSet = newCodeSet()

func registerCodec(codecs ...Codec) {
	codecSet.registerCodec(codecs...)
}

func (cs CodecSet) registerCodec(codecs ...Codec) {
	for _, codec := range codecs {
		cs[reflectTypeName(codec)] = codec
	}
}

func getCodec(name string) (Codec, error) {
	return codecSet.getCodec(name)
}

func (cs CodecSet) getCodec(name string) (Codec, error) {
	if codec, ok := cs[name]; ok {
		return codec, nil
	}
	return nil, ErrEncodingMissingCodec
}

func reflectTypeName(i interface{}) string {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
		return reflect.TypeOf(i).Elem().Name()
	} else if v.Kind() == reflect.Struct {
		return reflect.TypeOf(i).Name()
	}
	return ""
}

type Encoding struct {
	encoder []string
	decoder []string
}

func NewEncoding(opts ...encodingOpt) *Encoding {
	e := &Encoding{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e Encoding) Reverse() Encoding {
	re := Encoding{
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

func Encode(e interface{}, v interface{}) ([]byte, error) {
	bytes, err := Marshal(e, v)
	return bytes.Data, err
}

func Decode(e interface{}, data []byte, v interface{}) error {
	return Unmarshal(e, MakeBytes(data), v)
}

func Marshal(e interface{}, v interface{}) (Bytes, error) {
	switch e := e.(type) {
	case Encoding:
		return e.Marshal(v)
	case Codec:
		return e.Marshal(v)
	default:
		return nilBytes, ErrMarshallerUnknownType
	}
}

func Unmarshal(e interface{}, data Bytes, v interface{}) error {
	switch e := e.(type) {
	case Encoding:
		return e.Unmarshal(data, v)
	case Codec:
		return e.Unmarshal(data, v)
	default:
		return ErrUnmarshallerUnknownType
	}
}

func (e Encoding) Marshal(v interface{}) (Bytes, error) {
	var bytes Bytes
	for index, name := range e.encoder {
		codec, err := getCodec(name)
		if err != nil {
			return nilBytes, err
		}
		if index == 0 {
			bytes, err = codec.Marshal(v)
			if err != nil {
				return nilBytes, err
			}
		} else {
			bytes, err = codec.Marshal(bytes)
			if err != nil {
				return nilBytes, err
			}
		}

	}
	return bytes, nil
}

func (e Encoding) Unmarshal(bytes Bytes, v interface{}) error {
	for index, name := range e.decoder {
		codec, err := getCodec(name)
		if err != nil {
			return err
		}
		if index == len(e.decoder)-1 {
			err = codec.Unmarshal(bytes, v)
			if err != nil {
				return err
			}
		} else {
			err = codec.Unmarshal(bytes, &bytes)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type encodingOpt func(*Encoding)
type encodingOption struct{}

var EncodingOption encodingOption

func (encodingOption) WithEncoder(path string, codecSep magic.SeparatorType, wordSep magic.SeparatorType) encodingOpt {
	return func(e *Encoding) {
		e.WithEncoder(path, codecSep, wordSep)
	}
}

func (e *Encoding) WithEncoder(path string, codecSep magic.SeparatorType, wordSep magic.SeparatorType) *Encoding {
	names := strings.Split(path, codecSep)
	for _, name := range names {
		e.encoder = append(e.encoder, magic.Standardize(name, wordSep))
	}
	return e
}

func (encodingOption) WithDecoder(path string, codecSep magic.SeparatorType, wordSep magic.SeparatorType) encodingOpt {
	return func(e *Encoding) {
		e.WithDecoder(path, codecSep, wordSep)
	}
}

func (e *Encoding) WithDecoder(path string, codecSep magic.SeparatorType, wordSep magic.SeparatorType) *Encoding {
	names := strings.Split(path, codecSep)
	for _, name := range names {
		e.decoder = append(e.decoder, magic.Standardize(name, wordSep))
	}
	return e
}
