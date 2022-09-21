package hierarchy

import (
	"bytes"
	"path/filepath"
	"sort"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func ReadEnv(prefix string) error {
	return _default.ReadEnv(prefix)
}

func (h *Hierarchy) ReadEnv(prefix string) error {
	h.AutomaticEnv()
	h.SetEnvPrefix(prefix)

	return nil
}

func ReadFlags(flags *pflag.FlagSet) error {
	return _default.ReadFlags(flags)
}

func (h *Hierarchy) ReadFlags(flags *pflag.FlagSet) error {
	return BindPFlags(flags)
}

func ReadConfigMap(m map[string][]byte) error {
	return _default.ReadAssetMap(m)
}

func (h *Hierarchy) ReadAssetMap(configMap map[string][]byte) error {
	keys := make([]string, 0, len(configMap))
	for name := range configMap {
		keys = append(keys, name)
	}

	sort.Strings(keys)

	for _, name := range keys {
		ext := filepath.Ext(name)
		data := configMap[name]

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
