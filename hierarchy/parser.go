package hierarchy

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// examples:
//  - [txt](project:aries,runtime=default)
//     Set: project = aries, runtime = default
//  - [lua]({a=1,b=2})
//     Set: a = 1, b = 2
//  - [json]({"a":1,"b":2})
//     Set: a = 1, b = 2
//  - [boot]{args}
//     Boot from args
//  - [auto]{flags}
//     Read from flags
//  - [yaml]{stdin}
//     Read from stdin.
//  - [boot]{stdin}
//     Boot from stdin
//  - [auto]{env:Aries}
// 	   Read from env with prefix Aries.
//  - [json]{cluster:default}
//     Read from cluster with name default.
//  - [boot]{hierarchy:boot_file}
//     Boot from hierarchy key
//  - [json]<http://filestone.com/file.json>
// 	   Read from http.
//  - [ini]<file:///E:/Filename/file.ini> ...
//	   Read from local file.
//  - [yaml]<etcd://192.168.1.2:2379@usr:passwd/aries/hierarchy>
//     Read from etcd.
//  - [yml]<redis://192.168.1.2:6379@usr:passwd/0/aries/hierarchy>
//     Read from redis.

var (
	ErrArgNotMatch   = errors.New("arg not match")
	ErrInvalidString = errors.New("invalid string")
)

func IsArgNotMatch(err error) bool {
	return errors.Is(err, ErrArgNotMatch)
}

type Parser struct {
	hierarchy *Hierarchy
}

func (*Parser) JSON(data []byte) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("json")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	return v, nil
}

func (*Parser) YAML(data []byte) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	return v, nil
}

func (*Parser) YML(data []byte) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("yml")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	return v, nil
}

func (*Parser) TOML(data []byte) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("toml")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	return v, nil
}

func (*Parser) HCL(data []byte) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("hcl")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	return v, nil
}

func (*Parser) Properties(data []byte) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("properties")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	return v, nil
}

func (*Parser) Props(data []byte) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("props")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	return v, nil
}

func (*Parser) Prop(data []byte) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("prop")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	return v, nil
}

func (*Parser) TFVars(data []byte) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("tfvars")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	return v, nil
}

func (*Parser) DotEnv(data []byte) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("dotenv")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	return v, nil
}

func (*Parser) Env(data []byte) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("env")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	return v, nil
}

func (*Parser) INI(data []byte) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("ini")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	return v, nil
}

func (*Parser) TXT(data []byte) (*viper.Viper, error) {
	v := viper.New()
	s := string(data)

	strs := []string{}

	for _, s := range strings.Split(s, "\n") {
		for _, s := range strings.Split(s, ",") {
			strs = append(strs, strings.TrimSpace(s))
		}
	}

	for _, str := range strs {
		if str == "" {
			continue
		}

		kvStrs := strings.SplitN(str, "=", 2)
		if len(kvStrs) != 2 {
			kvStrs = strings.SplitN(str, ":", 2)
			if len(kvStrs) != 2 {
				return nil, fmt.Errorf("%w: %s", ErrInvalidString, s)
			}
		}

		v.Set(kvStrs[0], kvStrs[1])
	}

	return v, nil
}

func (*Parser) Lua(data []byte) (*viper.Viper, error) {
	panic("not implemented")
}

func (*Parser) JS(data []byte) (*viper.Viper, error) {
	panic("not implemented")
}

func (*Parser) Javascript(data []byte) (*viper.Viper, error) {
	panic("not implemented")
}

func (*Parser) Python(data []byte) (*viper.Viper, error) {
	panic("not implemented")
}

func (*Parser) Py(data []byte) (*viper.Viper, error) {
	panic("not implemented")
}

func (*Parser) Python3(data []byte) (*viper.Viper, error) {
	panic("not implemented")
}

func (*Parser) Py3(data []byte) (*viper.Viper, error) {
	panic("not implemented")
}

func (*Parser) Boot(data []byte) (*viper.Viper, error) {
	panic("not implemented")
}

func (p *Parser) Auto(data []byte) (*viper.Viper, error) {
	strs := strings.Split(string(data), ":")
	switch strs[0] {
	case "args":
	case "env":
		if len(strs) < 2 {
			return nil, fmt.Errorf("%w: %s", ErrModeAssetsArgsNotEnough, string(data))
		}
		p.hierarchy.AutomaticEnv()
		p.hierarchy.SetEnvPrefix(strs[1])
	case "flags":
		panic("not implemented")
	}
	return viper.New(), nil
}
