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

func (Base64) Marshal(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case []byte:
		s := base64.StdEncoding.EncodeToString(v)
		return []byte(s), nil
	case Bytes:
		s := base64.StdEncoding.EncodeToString(v.Data)
		return []byte(s), nil
	case *Bytes:
		s := base64.StdEncoding.EncodeToString(v.Data)
		return []byte(s), nil
	default:
		return nil, ErrBase64WrongValueType
	}
}

func (Base64) Unmarshal(data []byte, v interface{}) error {
	switch v := v.(type) {
	case *Bytes:
		s, err := base64.StdEncoding.DecodeString(string(data))
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

func init() {
	register(new(Base64URL))
}

func (Base64URL) Marshal(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case []byte:
		s := base64.URLEncoding.EncodeToString(v)
		return []byte(s), nil
	case Bytes:
		s := base64.URLEncoding.EncodeToString(v.Data)
		return []byte(s), nil
	case *Bytes:
		s := base64.URLEncoding.EncodeToString(v.Data)
		return []byte(s), nil
	default:
		return nil, ErrBase64URLWrongValueType
	}
}

func (Base64URL) Unmarshal(data []byte, v interface{}) error {
	switch v := v.(type) {
	case *Bytes:
		s, err := base64.URLEncoding.DecodeString(string(data))
		if err != nil {
			return err
		}
		v.Data = []byte(s)
		return nil
	default:
		return ErrBase64URLWrongValueType
	}
}
