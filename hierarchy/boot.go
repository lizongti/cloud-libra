package hierarchy

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/pflag"
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

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(config string) (*viper.Viper, error) {
	panic("implement me")
}

func ReadArgs(args []string) error {
	return _default.ReadArgs(args)
}

func (h *Hierarchy) ReadArgs(args []string) error {
	for _, arg := range args {
		v, err := NewParser().Parse(arg)
		if IsArgNotMatch(err) {
			continue
		} else if err != nil {
			return err
		}

		if err := h.MergeConfigMap(v.AllSettings()); err != nil {
			return err
		}
	}

	return nil
}

func ReadHierarchyValue(key string) error {
	return _default.ReadHierarchyValue(key)
}

func (h *Hierarchy) ReadHierarchyValue(key string) error {
	str := h.GetString(key)
	for _, arg := range strings.Split(str, " ") {
		v, err := NewParser().Parse(arg)
		if IsArgNotMatch(err) {
			continue
		} else if err != nil {
			return err
		}

		if err := h.MergeConfigMap(v.AllSettings()); err != nil {
			return err
		}
	}

	return nil
}

func (h *Hierarchy) ReadEnv(prefix string) error {
	h.AutomaticEnv()
	h.SetEnvPrefix(prefix)

	return nil
}

func ReadAssetMap(assetMap map[string][]byte) error {
	return _default.ReadAssetMap(assetMap)
}

func (h *Hierarchy) ReadAssetMap(assetMap map[string][]byte) error {
	keys := make([]string, 0, len(assetMap))
	for name := range assetMap {
		keys = append(keys, name)
	}

	sort.Strings(keys)

	for _, name := range keys {
		ext := filepath.Ext(name)
		data := assetMap[name]

		v := viper.New()
		v.SetConfigType(ext[1:])

		if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
			return err
		}

		if err := h.MergeConfigMap(v.AllSettings()); err != nil {
			return err
		}
	}

	return nil
}

func ReadFlags(flags *pflag.FlagSet) error {
	return _default.ReadFlags(flags)
}

func (h *Hierarchy) ReadFlags(flags *pflag.FlagSet) error {
	return BindPFlags(flags)
}

func ReadStdin(typ string) error {
	return _default.ReadStdin(typ)
}

func (h *Hierarchy) ReadStdin(typ string) error {
	v := viper.New()
	v.SetConfigType(typ)

	if data, err := io.ReadAll(bufio.NewReader(os.Stdin)); err != nil {
		return err
	} else if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return err
	}

	return h.MergeConfigMap(v.AllSettings())
}

func ReadPlainText(str string) error {
	return _default.ReadPlainText(str)
}

func (h *Hierarchy) ReadPlainText(str string) error {
	v := viper.New()

	for _, kvStr := range strings.Split(str, ",") {
		kvStrs := strings.SplitN(kvStr, "=", 2)
		if len(kvStrs) != 2 {
			kvStrs = strings.SplitN(kvStr, ":", 2)
			if len(kvStrs) != 2 {
				return fmt.Errorf("%w: %s", ErrInvalidString, kvStr)
			}
		}

		v.Set(kvStrs[0], kvStrs[1])
	}

	return h.MergeConfigMap(v.AllSettings())
}

func ReadEncodedText(typ string, str string) error {
	return _default.ReadEncodedText(typ, str)
}

func (h *Hierarchy) ReadEncodedText(typ string, str string) error {
	v := viper.New()
	v.SetConfigType(typ)

	if err := v.ReadConfig(bytes.NewReader([]byte(str))); err != nil {
		return err
	}

	return h.MergeConfigMap(v.AllSettings())
}

func ReadScriptText(typ string, script string) error {
	return _default.ReadScript(typ, script)
}

func (h *Hierarchy) ReadScript(typ, script string) error {
	panic("implement me")
}

func ReadURL(url string) error {
	return _default.ReadURL(url)
}

func (h *Hierarchy) ReadURL(url string) error {
	panic("implement me")
}

func ReadCluster(prefix string) error {
	return _default.ReadCluster(prefix)
}

func (h *Hierarchy) ReadCluster(prefix string) error {
	panic("implement me")
}
