package settings

import (
	"github.com/spf13/viper"
)

var (
	Env       = viper.New()
	RootFlags = viper.New()
	Flags     = viper.New()
)

func init() {
	Env.AutomaticEnv()

}
