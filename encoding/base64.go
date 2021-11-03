package encoding

import (
	"encoding/base64"
	"errors"
)

var (
	ErrBase64WrongValueType    = errors.New("codec base64 converts on wrong type value")
	ErrBase64URLWrongValueType = errors.New("codec base64URL converts on wrong type value")
)

type Base64 struct{}

func init() {
	register(new(Base64))
}

func (Base64) Marshal(v interface{}) (Bytes, error) {
	switch v := v.(type) {
	case Bytes:
		s := base64.StdEncoding.EncodeToString(v.Data)
		bytes := MakeBytes([]byte(s))
		return bytes, nil
	case *Bytes:
		s := base64.StdEncoding.EncodeToString(v.Data)
		bytes := MakeBytes([]byte(s))
		return bytes, nil
	default:
		return nilBytes, ErrBase64WrongValueType
	}
}

func (Base64) Unmarshal(data Bytes, v interface{}) error {
	switch v := v.(type) {
	case *Bytes:
		s, err := base64.StdEncoding.DecodeString(string(v.Data))
		if err != nil {
			return err
		}
		v.Data = []byte(s)
		return nil
	default:
		return ErrBase64WrongValueType
	}
}

type Base64URL struct{}

var base64URL = new(Base64URL)

func (Base64URL) String() string {
	return "base64url"
}

func (Base64URL) Marshal(v interface{}) (Bytes, error) {
	switch v := v.(type) {
	case Bytes:
		s := base64.URLEncoding.EncodeToString(v.Data)
		bytes := MakeBytes([]byte(s))
		return bytes, nil
	case *Bytes:
		s := base64.URLEncoding.EncodeToString(v.Data)
		bytes := MakeBytes([]byte(s))
		return bytes, nil
	default:
		return nilBytes, ErrBase64URLWrongValueType
	}
}

func (Base64URL) Unmarshal(data Bytes, v interface{}) error {
	switch v := v.(type) {
	case *Bytes:
		data, err := base64.URLEncoding.DecodeString(string(v.Data))
		if err != nil {
			return err
		}
		v.Data = data
		return nil
	default:
		return ErrBase64URLWrongValueType
	}
}
