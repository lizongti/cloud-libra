package global

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"

	"github.com/cloudlibraries/libra/osutil"
	"github.com/spf13/cobra"
)

/*
Viper uses the following precedence order. Each item takes precedence over the item below it:
* explicit call to Set
* flag
* env
* config
* key/value store
* default
Important: Viper configuration keys are case insensitive. There are ongoing discussions about making that optional.
*/

func init() {
	initVars()
	initEnv()
	initConfig()
	initRemoteConfig()
}

func initVars() {
	Set(Project, "aries")
	Set(ProjectPath, osutil.GetProjectPath(GetString(Project)))
	Set(WorkPath, osutil.GetWorkPath())
	Set(ExePath, osutil.GetExePath())
	Set(ConfigRelPath, filepath.Join("assets", "config"))
	Set(ConfigPath, filepath.Join(GetString(ProjectPath), GetString(ConfigRelPath)))
	Set(ConfigType, "yml")
}

func initEnv() {
	AutomaticEnv()
	SetEnvPrefix(fmt.Sprintf("%s.", GetString(Project)))
}

func BindCmd(cmd *cobra.Command) {
	if err := BindPFlags(cmd.Flags()); err != nil {
		log.Panic(err)
	}
}

func initConfig() {
	configType := GetString(ConfigType)
	configRelPath := GetString(ConfigRelPath)
	configPath := GetString(ConfigPath)
	runtime := GetString(Runtime)

	SetConfigType(configType)
	stat := osutil.Stat(configPath)
	switch stat.Type() {
	case osutil.StatDir:
		for _, data := range osutil.ReadFilesBySuffix(configType,
			filepath.Join(configPath, "generic"),
			filepath.Join(configPath, "runtime", runtime),
		) {
			if err := MergeConfig(bytes.NewReader(data)); err != nil {
				log.Panic(err)
			}
		}
	case osutil.StatNotExists:
		assetNames := AssetNames()
		for _, assetName := range assetNames {
			if osutil.IsParent(assetName, filepath.Join(configRelPath, "generic")) ||
				osutil.IsParent(assetName, filepath.Join(configRelPath, "runtime", runtime)) {
				data, err := Asset(assetName)
				if err != nil {
					log.Panic(err)
				}
				if err := MergeConfig(bytes.NewReader(data)); err != nil {
					log.Panic(err)
				}
			}
		}
	default:
		log.Panic(stat.Err())
	}
}

// initRemoteConfig initializes remote config from etcdv3.
func initRemoteConfig() {
	configType := GetString(ConfigType)
	SetConfigType(configType)

	if err := AddRemoteProvider("etcd3", GetString(RemoteEndpoint), "/test"); err != nil {
		log.Panic(err)
	}
	if err := ReadRemoteConfig(); err != nil {
		log.Panic(err)
	}
}
