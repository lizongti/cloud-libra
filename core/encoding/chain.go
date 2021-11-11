package encoding

import (
	"strings"

	"github.com/aceaura/libra/magic"
)

type Chain struct {
	encoder []string
	decoder []string
}

func NewChain(opts ...funcChainOption) *Chain {
	e := &Chain{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

var empty *Chain = NewChain()

func Empty() *Chain {
	return empty
}

func (c Chain) String() string {
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

func (c Chain) Reverse() Encoding {
	re := Chain{
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

func (c Chain) Marshal(v interface{}) ([]byte, error) {
	var data []byte
	for index, name := range c.encoder {
		encoding, err := localEncoding(name)
		if err != nil {
			return nil, err
		}
		if index == 0 {
			data, err = encoding.Marshal(v)
			if err != nil {
				return nil, err
			}
		} else {
			data, err = encoding.Marshal(data)
			if err != nil {
				return nil, err
			}
		}

	}
	return data, nil
}

func (c Chain) Unmarshal(data []byte, v interface{}) error {
	bytes := MakeBytes(nil)
	for index, name := range c.decoder {
		encoding, err := localEncoding(name)
		if err != nil {
			return err
		}
		if index < len(c.decoder)-1 {
			err = encoding.Unmarshal(data, &bytes)
			if err != nil {
				return err
			}
			data = bytes.Data
		} else {
			err = encoding.Unmarshal(data, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type funcChainOption func(*Chain)
type chainOption struct{}

var ChainOption chainOption

func (chainOption) WithEncoder(path string, encodingSep magic.SeparatorType, wordSep magic.SeparatorType) funcChainOption {
	return func(c *Chain) {
		c.WithEncoder(path, encodingSep, wordSep)
	}
}

func (c *Chain) WithEncoder(path string, encodingSep magic.SeparatorType, wordSep magic.SeparatorType) *Chain {
	names := strings.Split(path, encodingSep)
	for _, name := range names {
		c.encoder = append(c.encoder, magic.Standardize(name, wordSep))
	}
	return c
}

func (chainOption) WithDecoder(path string, encodingSep magic.SeparatorType, wordSep magic.SeparatorType) funcChainOption {
	return func(c *Chain) {
		c.WithDecoder(path, encodingSep, wordSep)
	}
}

func (c *Chain) WithDecoder(path string, encodingSep magic.SeparatorType, wordSep magic.SeparatorType) *Chain {
	names := strings.Split(path, encodingSep)
	for _, name := range names {
		c.decoder = append(c.decoder, magic.Standardize(name, wordSep))
	}
	return c
}
