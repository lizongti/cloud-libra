package hierarchy

import (
	"bytes"
	"path/filepath"
	"sort"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func LoadEnv(prefix string) error {
	return _default.LoadEnv(prefix)
}

func (h *Hierarchy) LoadEnv(prefix string) error {
	h.AutomaticEnv()
	h.SetEnvPrefix(prefix)

	return nil
}

func LoadFlags(flags *pflag.FlagSet) error {
	return _default.LoadFlags(flags)
}

func (h *Hierarchy) LoadFlags(flags *pflag.FlagSet) error {
	return BindPFlags(flags)
}

func LoadConfigMap(m map[string][]byte) error {
	return _default.LoadAssetMap(m)
}

func (h *Hierarchy) LoadAssetMap(assetMap map[string][]byte) error {
	keys := make([]string, 0, len(assetMap))
	for name := range assetMap {
		keys = append(keys, name)
	}

	sort.Strings(keys)

	for _, name := range keys {
		ext := filepath.Ext(name)
		data := h.ReplaceAllVars(assetMap[name])

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
