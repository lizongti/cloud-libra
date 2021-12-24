package encoding

import (
	"gopkg.in/yaml.v2"

	"github.com/aceaura/libra/core/magic"
)

type YAML struct{}

func init() {
	register(new(YAML))
}

func NewYAML() *YAML {
	return new(YAML)
}

func (yaml YAML) String() string {
	return magic.TypeName(yaml)
}

func (yaml YAML) Style() EncodingStyleType {
	return EncodingStyleStruct
}

func (YAML) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (YAML) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

func (yaml YAML) Reverse() Encoding {
	return yaml
}
