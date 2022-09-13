package settings

import (
	"io"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Configuration struct {
	*viper.Viper
}

/*
explicit call to Set
flag
env
config
key/value store
default
*/

func NewConfiguration(typ string, readers []io.Reader) *Configuration {
	c := &Configuration{
		Viper: viper.New(),
	}

	// Config
	c.SetConfigType(typ)
	for _, reader := range readers {
		if err := c.MergeConfig(reader); err != nil {
			log.Panic(err)
		}
	}

	// Environment variables
	c.AutomaticEnv()

	// Flags
	serverCmd := &cobra.Command{}
	viper.BindPFlags(serverCmd.Flags())
	return c
}

// func (c *Configuration) initConfigFromFile() {
// 	c.SetConfigType("yaml")

// 	c.AddConfigPath(c.path)
// 	paths := WalkFilesBySuffix(c.path, ".yaml")
// 	relPaths := make([]string, 0, len(paths))
// 	for _, path := range paths {
// 		relPath := strings.TrimPrefix(path, c.path)
// 		relPaths = append(relPaths, relPath)
// 	}

// 	for _, relPath := range relPaths {
// 		if strings.HasPrefix(relPath, "generic") {
// 			c.SetConfigName(relPath)
// 			if err := c.MergeInConfig(); err != nil {
// 				panic(err)
// 			}
// 		}
// 	}

// 	for _, relPath := range relPaths {
// 		if strings.HasPrefix(relPath, filepath.Join("runtime", c.runtime)) {
// 			c.SetConfigName(relPath)
// 			if err := c.MergeInConfig(); err != nil {
// 				panic(err)
// 			}
// 		}
// 	}
// }
